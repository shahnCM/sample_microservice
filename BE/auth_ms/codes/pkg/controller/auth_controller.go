package controller

import (
	"auth_ms/pkg/action"
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/helper/common"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FreshToken(ctx *fiber.Ctx) error {

	// Get Post Body
	userLoginReqP := new(request.UserLoginDto)
	if errBody, err := common.ParseRequestBody(ctx, userLoginReqP); errBody != nil {
		return err
	}

	tokenDataP, err := action.Login(userLoginReqP)
	if err != nil {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, err.Error())
	}

	// return JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, tokenDataP, nil, nil)
}

func RegisterUser(ctx *fiber.Ctx) error {

	// Get Post Body
	userP := new(request.UserRegistrationDto)
	if errBody, err := common.ParseRequestBody(ctx, userP); errBody != nil {
		return err
	}
	// Check password with confirm password
	if userP.Password != userP.PasswordConfirm {
		log.Println(userP.Password != userP.PasswordConfirm)
		return common.ErrorResponse(ctx, fiber.ErrUnprocessableEntity.Code, "Password mismatch")
	}

	if err := action.Register(userP); err != nil {
		return common.ErrorResponse(ctx, err.Code, err.Message)
	}

	return common.SuccessResponse(ctx, 201, nil, nil, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {

	// Jwt Token from POST body
	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, "Invalid Jwt Token")
	}

	// Refresh Token from POST body
	refreshTokenReqP := new(request.RefreshTokenDto)
	if errBody, validateionErr := common.ParseRequestBody(ctx, refreshTokenReqP); errBody != nil {
		return validateionErr
	}
	refreshToken := *refreshTokenReqP.Token

	tokenDataP, err := action.Refresh(&jwtToken, &refreshToken)
	if err != nil {
		return common.ErrorResponse(ctx, err.Code, err.Message)
	}

	// return renewed JWT and Refresh Token
	return common.SuccessResponse(ctx, 200, tokenDataP, nil, nil)
}

func VerifyToken(ctx *fiber.Ctx) error {
	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	userTokenDataP, err := action.Verify(&jwtToken)
	if err != nil {
		return common.ErrorResponse(ctx, err.Code, err.Message)
	}

	return common.SuccessResponse(ctx, 200, userTokenDataP, nil, nil)
}

func RevokeToken(ctx *fiber.Ctx) error {
	jwtToken, _ := common.ParseHeader(ctx, "Authorization", "Bearer ")
	jwtToken = strings.Split(jwtToken, "Bearer ")[1]
	if jwtToken == "" {
		return common.ErrorResponse(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message)
	}

	if err := action.Revoke(&jwtToken); err != nil {
		return common.ErrorResponse(ctx, err.Code, err.Message)
	}

	// returns 200 ok
	return common.SuccessResponse(ctx, 204, nil, nil, nil)
}
