package controller

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/model"
	"auth_ms/pkg/service"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func Login(ctx *fiber.Ctx) error {
	// Get Post Body
	userRegReqP := new(request.UserLoginDto)
	if errBody, err := common.ParseRequestBody(ctx, userRegReqP); errBody != nil {
		return err
	}

	// declaring service response and error
	var responseP *response.GenericServiceResponseDto
	var err error

	// Verify Username & Password
	responseP, err = service.GetUser(userRegReqP)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Credentials")
	}
	userModelP := responseP.Data.(*model.User)

	// Generate ULID for Token
	ulidP, err := common.GenerateULID()
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error")
	}

	// Set claims & Generate a new JWT token and Associated Refresh Token
	responseP, err = service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, ulidP)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error")
	}
	tokenDataP := responseP.Data.(*dto.TokenDataDto)

	// Asynchronously manage associated session & token
	defer safeasync.Run(func() {
		actionFailed := 0

		// Create a New Associated Session
		responseP, err = service.StoreSession(&userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
			actionFailed++
		}
		sessionModelP := responseP.Data.(*model.Session)

		// Create associated token
		responseP, err = service.StoreToken(&userModelP.Id, &sessionModelP.Id, ulidP, tokenDataP)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
			actionFailed++
		}
		tokenModelP := responseP.Data.(*model.Token)

		// Revoke last active token
		if userModelP.LastTokenId != nil {
			_, err = service.UpdateTokenStatus(userModelP.LastTokenId, enum.REVOCKED_TOKEN)
			if err != nil {
				log.Println("! CRITICAL " + err.Error())
				actionFailed++
			}
		}

		// End last active session
		if userModelP.LastSessionId != nil {
			_, err = service.EndSession(userModelP.LastSessionId)
			if err != nil {
				log.Println("! CRITICAL " + err.Error())
				actionFailed++
			}
		}

		// Update user with new session id and new token id
		_, err = service.UpdateUserActiveSessionAndToken(&userModelP.Id, &sessionModelP.Id, tokenModelP.Id)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
			actionFailed++
		}

		log.Println("Total failed action: ", actionFailed)
	})

	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, tokenDataP, nil, nil)
}

func Register(ctx *fiber.Ctx) error {
	// Get Post Body
	userP := new(request.UserRegistrationDto)
	if errBody, err := common.ParseRequestBody(ctx, userP); errBody != nil {
		return err
	}

	// Store Username & Password
	_, err := service.StoreUser(userP)
	if err != nil {
		// return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, strings.Split(err.Error(), ":")[1])
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, err.Error())
	}

	return common.SuccessResponse(ctx, 201, nil, nil, nil)
}

func Refresh(ctx *fiber.Ctx) error {

	// Jwt Token from POST body
	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	// Refresh Token from POST body
	refreshTokenReqP := new(dto.RefreshTokenDto)
	if errBody, err := common.ParseRequestBody(ctx, refreshTokenReqP); errBody != nil {
		return err
	}
	refreshToken := *refreshTokenReqP.Token

	// declaring service response and error
	var responseP *response.GenericServiceResponseDto
	// var err error
	var wg sync.WaitGroup
	wg.Add(2)
	statusResults := make(chan int, 2)
	claimsResults := make(chan *service.Claims, 2)

	// verify jwt and refresh token concurrently
	safeasync.Run(func() {
		defer wg.Done()
		responseP, _ = service.VerifyJWT(&jwtToken)
		if responseP != nil {
			statusResults <- (responseP.StatusCode - 1)
			claimsResults <- responseP.Data.(*service.Claims)
		} else {
			statusResults <- 0
			claimsResults <- &service.Claims{}
		}
	})

	safeasync.Run(func() {
		defer wg.Done()
		responseP, _ = service.VerifyRefreshToken(&refreshToken)
		if responseP != nil {
			statusResults <- (responseP.StatusCode + 9)
			claimsResults <- responseP.Data.(*service.Claims)
		} else {
			statusResults <- 0
			claimsResults <- &service.Claims{}
		}
	})

	wg.Wait()
	close(statusResults)
	close(claimsResults)

	totalStatusResults := 0
	for result := range statusResults {
		totalStatusResults += result
	}

	var claimsArr []*service.Claims
	for result := range claimsResults {
		claimsArr = append(claimsArr, result)
	}
	claimsOk := *claimsArr[0].TokenId == *claimsArr[1].TokenId

	if !claimsOk || totalStatusResults != 609 {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid refresh token")
	}

	/**
	 *check database fix session or token faults if there's any
	 */
	return common.SuccessResponse(ctx, 200, "auth - Refresh", nil, nil)
}

func Verify(ctx *fiber.Ctx) error {
	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	// declaring service response and error
	var responseP *response.GenericServiceResponseDto
	var err error

	// Verify token and get claims
	responseP, err = service.VerifyJWT(&jwtToken)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}
	if responseP.StatusCode != 200 {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Token expired, please refresh")
	}
	claims := responseP.Data.(*service.Claims)

	// Check if the token ulid exists on db and its status is fresh
	responseP, err = service.GetToken(claims.TokenId)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, err.Error())
	}
	tokenModelP := responseP.Data.(*model.Token)
	if tokenModelP.TokenStatus != enum.FRESH_TOKEN { // Allow fresh only
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Token: Not fresh")
	}

	/**
	 * Now check db (users table) if its an active token
	 * else revoke this token
	 * And associated session
	 * return unauthorized
	 */

	// Fetch user by user_id
	responseP, err = service.GetUserById(&claims.UserId)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Token")
	}
	userModelP := responseP.Data.(*model.User)

	// Check if user's active session_id & token_id matches with the claim
	if *userModelP.LastSessionId != *tokenModelP.SessionId || *userModelP.LastTokenId != *tokenModelP.Id {

		// if doesn't match revoke token and end session
		defer safeasync.Run(func() {
			actionFailed := 0

			// Revoke requested token
			_, err = service.UpdateTokenStatus(tokenModelP.Id, enum.REVOCKED_TOKEN)
			if err != nil {
				log.Println("! CRITICAL " + err.Error())
				actionFailed++
			}

			// End session associated with requested token

			_, err = service.EndSession(tokenModelP.SessionId)
			if err != nil {
				log.Println("! CRITICAL " + err.Error())
				actionFailed++
			}

			log.Println("Total failed action: ", actionFailed)
		})

		return common.ErrorResponse(ctx, 401, fiber.ErrUnauthorized.Message)
	}

	return common.SuccessResponse(ctx, 200, nil, nil, nil)
}

func Logout(ctx *fiber.Ctx) error {
	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	// Declaring service response and error
	var responseP *response.GenericServiceResponseDto
	var err error

	// Verify token and get claims
	responseP, err = service.VerifyJWT(&jwtToken)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}
	claims := responseP.Data.(*service.Claims)

	defer safeasync.Run(func() {
		actionFailed := 0

		// Fetch user by user_id
		responseP, err = service.GetUserById(&claims.UserId)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
		}
		userModelP := responseP.Data.(*model.User)

		// Check if the token ulid exists on db
		responseP, err = service.GetToken(claims.TokenId)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
		}
		tokenModelP := responseP.Data.(*model.Token)

		// Revoke token if not already revoked
		if tokenModelP.TokenStatus != enum.REVOCKED_TOKEN {
			_, err = service.UpdateTokenStatus(tokenModelP.Id, enum.REVOCKED_TOKEN)
			if err != nil {
				log.Println("! CRITICAL " + err.Error())
			}
		}
		// Update user's active session to NULL Only if its associated with the requested token
		if *userModelP.LastSessionId == *tokenModelP.SessionId {
			_, err = service.UpdateUserActiveSessionAndToken(tokenModelP.UserId, nil, nil)
			if err != nil {
				log.Println("! CRITICAL " + err.Error())
			}
		}
		// End session associated with the token
		_, err = service.EndSession(tokenModelP.SessionId)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
		}

		log.Println("Total failed action: ", actionFailed)
	})

	return common.SuccessResponse(ctx, 204, nil, nil, nil)
}
