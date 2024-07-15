package controller

import (
	"auth_ms/pkg/dto"
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/helper/common"
	"auth_ms/pkg/helper/safeasync"
	"auth_ms/pkg/model"
	"auth_ms/pkg/service"
	"log"

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
		// Create a New Associated Session
		responseP, err = service.StoreSession(&userModelP.Id)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
			return
		}
		// create associated token
		sessionP := responseP.Data.(*model.Session)
		responseP, err = service.StoreToken(&userModelP.Id, &sessionP.Id, ulidP, tokenDataP)
		if err != nil {
			log.Println("! CRITICAL " + err.Error())
			return
		}
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

	// Verify Username & Password
	_, err := service.StoreUser(userP)
	if err != nil {
		// return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, strings.Split(err.Error(), ":")[1])
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, err.Error())
	}

	return common.SuccessResponse(ctx, 201, nil, nil, nil)
}

func Refresh(ctx *fiber.Ctx) error {
	// Verify token and get claims
	// Check if the token ulid exists on db and its status is fresh
	// Update this token's status to refreshed
	// Generate a new JWT token and Associated Refresh Token
	// Associate it with the current Running session
	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, "auth - Refresh", nil, nil)
}

func Verify(ctx *fiber.Ctx) error {
	// Verify token and get claims
	// Check if the token ulid exists on db and its status is fresh
	return common.SuccessResponse(ctx, 200, "auth - Verify", nil, nil)
}

func Revoke(ctx *fiber.Ctx) error {
	// Verify token and get claims
	// Check if the token ulid exists on db and its status is fresh
	// Update this token's status to refreshed
	// Generate a new JWT token and Associated Refresh Token
	// Associate it with the current Running session
	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, "auth - Revoke", nil, nil)
}
