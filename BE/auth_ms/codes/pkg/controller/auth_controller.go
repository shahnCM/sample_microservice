package controller

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"auth_ms/pkg/service"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func FreshToken(ctx *fiber.Ctx) error {
	// Get Post Body
	userRegReqP := new(request.UserLoginDto)
	if errBody, err := common.ParseRequestBody(ctx, userRegReqP); errBody != nil {
		return err
	}

	// declaring service response and error
	var responseP *response.GenericServiceResponseDto
	var err error

	// Verify Username & Password
	responseP, err = service.GetUser(nil, userRegReqP)
	if err != nil {
		log.Println(err)
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

		tx := mariadb10.GetMariaDb10().Begin()
		if tx.Error != nil {
			log.Println("! CRITICAL Could not start transaction:", tx.Error)
			return
		}

		// Create a New Associated Session
		responseP, err = service.StoreSession(tx, &userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp)
		if err != nil {
			tx.Rollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}
		sessionModelP := responseP.Data.(*model.Session)

		// Update user with new session id and new token id
		_, err = service.UpdateUserActiveSessionAndToken(tx, &userModelP.Id, &sessionModelP.Id, ulidP)
		if err != nil {
			tx.Rollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// End last active session
		if userModelP.LastSessionId != nil {
			_, err = service.EndSession(tx, userModelP.LastSessionId)
			if err != nil {
				tx.Rollback()
				log.Println("! CRITICAL " + err.Error())
				return
			}
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			log.Println("! CRITICAL Transaction commit failed:", err.Error())
			return
		}
	})

	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, tokenDataP, nil, nil)
}

func RegisterUser(ctx *fiber.Ctx) error {
	// Get Post Body
	userP := new(request.UserRegistrationDto)
	if errBody, err := common.ParseRequestBody(ctx, userP); errBody != nil {
		return err
	}

	// Store Username & Password
	_, err := service.StoreUser(nil, userP)
	if err != nil {
		// return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, strings.Split(err.Error(), ":")[1])
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, err.Error())
	}

	return common.SuccessResponse(ctx, 201, nil, nil, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {

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
	var err error
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
			statusResults <- -1
			claimsResults <- &service.Claims{}
		}
	})

	// verify jwt and refresh token concurrently
	safeasync.Run(func() {
		defer wg.Done()
		responseP, _ = service.VerifyRefreshToken(&refreshToken)
		if responseP != nil {
			statusResults <- (responseP.StatusCode + 9)
			claimsResults <- responseP.Data.(*service.Claims)
		} else {
			statusResults <- -1
			claimsResults <- &service.Claims{}
		}
	})

	wg.Wait()
	close(statusResults)
	close(claimsResults)

	totalStatusResults := 0
	var claimsArr []*service.Claims
	for result := range statusResults {
		totalStatusResults += result
	}
	for result := range claimsResults {
		claimsArr = append(claimsArr, result)
	}
	claimsOk := *claimsArr[0].TokenId == *claimsArr[1].TokenId
	if !claimsOk || totalStatusResults != 609 {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid refresh token, 609")
	}
	claims := claimsArr[0]

	// check database fix session or token faults if there's any
	responseP, err = service.GetUserById(nil, &claims.UserId)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Invalid refresh token, No user with this token")
	}
	userModelP := responseP.Data.(*model.User)
	userLastSessionModelP := userModelP.LastSession

	// Check if user's active token_id matches with the claim's token_id
	if userModelP.LastTokenId == nil || *userModelP.LastTokenId != *claims.TokenId {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid refresh token, Inactive token")
	}

	// Generate ULID for Token
	ulidP, err := common.GenerateULID()
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error, ULID")
	}

	// Set claims & Generate a new JWT token and Associated Refresh Token
	responseP, err = service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, ulidP)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error, TOKEN ISSUE")
	}
	tokenDataP := responseP.Data.(*dto.TokenDataDto)

	// Asynchronously manage associated session & token
	defer safeasync.Run(func() {

		tx := mariadb10.GetMariaDb10().Begin()
		if tx.Error != nil {
			log.Println("! CRITICAL Could not start transaction:", tx.Error)
		}

		// Update active Associated Session's tokenId, jwtExpiresAt, refreshExpiresAt, refreshCount
		_, err = service.RefreshSession(tx, userModelP.LastSessionId, &userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp, &userLastSessionModelP.RefreshCount)
		if err != nil {
			tx.Rollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// Update user with new token id
		_, err = service.UpdateUserActiveToken(tx, &userModelP.Id, ulidP)
		if err != nil {
			tx.Rollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			log.Println("! CRITICAL Transaction commit failed:", err.Error())
			return
		}

	})

	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, tokenDataP, nil, nil)
}

func VerifyToken(ctx *fiber.Ctx) error {
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
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Token invalid, malformed")
	}
	if responseP.StatusCode != 200 {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Token invalid, expired")
	}
	claims := responseP.Data.(*service.Claims)

	// Fetch user by user_id from claims
	responseP, err = service.GetUserById(nil, &claims.UserId)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Token, foreign")
	}
	userModelP := responseP.Data.(*model.User)

	// Check if user's active token_id matches with the claim's token_id
	if userModelP.LastTokenId == nil || *userModelP.LastTokenId != *claims.TokenId {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	return common.SuccessResponse(ctx, 200, nil, nil, nil)
}

func RevokeToken(ctx *fiber.Ctx) error {
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

		// Fetch user by user_id
		responseP, err = service.GetUserById(nil, &claims.UserId)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
		}
		userModelP := responseP.Data.(*model.User)

		if userModelP.LastTokenId == nil || *userModelP.LastTokenId != *claims.TokenId {
			return
		}

		tx := mariadb10.GetMariaDb10().Begin()
		if tx.Error != nil {
			log.Println("! CRITICAL Could not start transaction:", tx.Error)
			return
		}

		// Update user's active session to NULL Only if its associated with the requested token
		_, err = service.UpdateUserActiveSessionAndToken(tx, &userModelP.Id, nil, nil)
		if err != nil {
			tx.Rollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// End session associated with the token
		_, err = service.EndSession(tx, userModelP.LastSessionId)
		if err != nil {
			tx.Rollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			log.Println("! CRITICAL Transaction commit failed:", err.Error())
			return
		}

	})

	return common.SuccessResponse(ctx, 204, nil, nil, nil)
}
