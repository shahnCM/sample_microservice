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
	userLoginReqP := new(request.UserLoginDto)
	if errBody, err := common.ParseRequestBody(ctx, userLoginReqP); errBody != nil {
		return err
	}

	// declaring service response and error
	var responseP *response.GenericServiceResponseDto
	var err error

	// Verify Username & Password
	responseP, err = service.GetUser(userLoginReqP)
	if err != nil {
		log.Println(err)
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Credentials")
	}
	userModelP := responseP.Data.(*model.User)

	// Compare password
	if !common.CompareHash(&userModelP.Password, &userLoginReqP.Password) {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Credentials")
	}

	// Generate ULID for Token
	ulidP, err := common.GenerateULID()
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error")
	}
	hashedUlid, err := common.GenerateHash(ulidP)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error")
	}

	// Set claims & Generate a new JWT token and Associated Refresh Token
	responseP, err = service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, hashedUlid)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error")
	}
	tokenDataP := responseP.Data.(*dto.TokenDataDto)

	// Asynchronously manage associated session & token
	defer safeasync.Run(func() {

		err := mariadb10.TransactionBegin().Error
		if err != nil {
			log.Println("! CRITICAL Could not start transaction:" + err.Error())
			return
		}

		// Create a New Associated Session
		responseP, err = service.StoreSession(&userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp)
		if err != nil {
			mariadb10.TransactionRollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}
		sessionModelP := responseP.Data.(*model.Session)

		// Update user with new session id and new token id
		_, err = service.UpdateUserActiveSessionAndToken(&userModelP.Id, &sessionModelP.Id, ulidP)
		if err != nil {
			mariadb10.TransactionRollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// End last active session
		if userModelP.LastSessionId != nil {
			_, err = service.EndSession(userModelP.LastSessionId)
			if err != nil {
				mariadb10.TransactionRollback()
				log.Println("! CRITICAL " + err.Error())
				return
			}
		}

		// Commit the transaction
		if err := mariadb10.TransactionCommit().Error; err != nil {
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

	if userP.Password != userP.PasswordConfirm {
		log.Println(userP.Password != userP.PasswordConfirm)
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, "Password mismatch")
	}
	hashedPassword, err := common.GenerateHash(&userP.Password)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message)
	}
	userP.Password = *hashedPassword

	// Store Username & Password
	_, err = service.StoreUser(userP)
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
	for result := range statusResults {
		totalStatusResults += result
	}
	if totalStatusResults != 609 {
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, "Can't refresh, Jwt isn't expired")
	}

	var claimsArr []*service.Claims
	for result := range claimsResults {
		claimsArr = append(claimsArr, result)
	}
	if *claimsArr[0].TokenId != *claimsArr[1].TokenId {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Refresh/Jwt token")
	}
	claims := claimsArr[0]

	// check database fix session or token faults if there's any
	// *** Lock users row here ***
	responseP, err = service.GetUserById(&claims.UserId)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Refresh token, User not found")
	}
	userModelP := responseP.Data.(*model.User)

	// Check if user's active token_id exists as there's any active session running
	if userModelP.LastTokenId == nil {
		// *** Release Lock ***
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, "No running session to refresh")
	}

	// Check if user's active token_id matches with the claim's token_id
	if !common.CompareHash(claims.TokenId, userModelP.LastTokenId) {
		// *** Release Lock ***
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Refresh token")
	}

	// Generate ULID for Token
	ulidP, err := common.GenerateULID()
	if err != nil {
		// *** Release Lock ***
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error, Id Error")
	}
	hashedUlid, err := common.GenerateHash(ulidP)
	if err != nil {
		// *** Release Lock ***
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error, Hash Error")
	}

	// Set claims & Generate a new JWT token and Associated Refresh Token
	responseP, err = service.IssueJwtWithRefreshToken(userModelP.Id, userModelP.Role, hashedUlid)
	if err != nil {
		// *** Release Lock ***
		return common.ErrorResponse(ctx, fiber.ErrInternalServerError.Code, "Internal Server Error, TOKEN ISSUE")
	}
	tokenDataP := responseP.Data.(*dto.TokenDataDto)

	// Asynchronously manage associated session & token
	defer safeasync.Run(func() {

		err = mariadb10.TransactionBegin().Error
		if err != nil {
			// *** Release Lock ***
			log.Println("! CRITICAL Could not start transaction:", err.Error())
		}

		// Update active Associated Session's tokenId, jwtExpiresAt, refreshExpiresAt, refreshCount
		_, err = service.RefreshSession(userModelP.LastSessionId, &userModelP.Id, ulidP, tokenDataP.Jwt.TokenExp, tokenDataP.Refresh.TokenExp, &userModelP.LastSession.RefreshCount)
		if err != nil {
			// *** Release Lock ***
			mariadb10.TransactionRollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// Update user with new token id
		_, err = service.UpdateUserActiveToken(&userModelP.Id, ulidP)
		if err != nil {
			// *** Release Lock ***
			mariadb10.TransactionRollback()
			log.Println("! CRITICAL " + err.Error())
			return
		}

		// Commit the transaction
		if err := mariadb10.TransactionCommit().Error; err != nil {
			// *** Release Lock ***
			log.Println("! CRITICAL Transaction commit failed:", err.Error())
			return
		}
		// *** Release Lock ***
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
	responseP, err = service.GetUserById(&claims.UserId)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Token, foreign")
	}
	userModelP := responseP.Data.(*model.User)

	// Check if user's active token_id matches with the claim's token_id
	if userModelP.LastTokenId == nil || !common.CompareHash(claims.TokenId, userModelP.LastTokenId) {
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
		responseP, err = service.GetUserById(&claims.UserId)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
		}
		userModelP := responseP.Data.(*model.User)

		if userModelP.LastTokenId == nil || !common.CompareHash(claims.TokenId, userModelP.LastTokenId) {
			return
		}

		err = mariadb10.TransactionBegin().Error
		if err != nil {
			log.Println("! CRITICAL Could not start transaction:", err.Error())
			return
		}

		// Update user's active session to NULL Only if its associated with the requested token
		_, err = service.UpdateUserActiveSessionAndToken(&userModelP.Id, nil, nil)
		if err != nil {
			mariadb10.TransactionRollback()
			log.Println("! CRITICAL ", err.Error())
			return
		}

		// End session associated with the token
		_, err = service.EndSession(userModelP.LastSessionId)
		if err != nil {
			mariadb10.TransactionRollback()
			log.Println("! CRITICAL ", err.Error())
			return
		}

		// Commit the transaction
		if err := mariadb10.TransactionCommit().Error; err != nil {
			log.Println("! CRITICAL Transaction commit failed:", err.Error())
			return
		}

	})

	return common.SuccessResponse(ctx, 204, nil, nil, nil)
}
