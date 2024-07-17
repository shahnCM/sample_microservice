package service

import (
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"

	"github.com/gofiber/fiber/v2"
)

func GetUserById(userIdP *uint) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindUserById(userIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: user}, nil
}

func GetUser(userP *request.UserLoginDto) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindUser(userP.Username, userP.Password)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: user}, nil
}

func StoreUser(userP *request.UserRegistrationDto) (*response.GenericServiceResponseDto, error) {
	if userP.Password != userP.PasswordConfirm {
		return nil, fiber.NewError(fiber.ErrUnprocessableEntity.Code, "Password mismatch")
	}

	userModelP := &model.User{
		Username:      userP.Username,
		Email:         userP.Email,
		Password:      userP.Password,
		Role:          enum.STANDARD,
		LastSessionId: nil,
		LastTokenId:   nil,
	}

	userRepo := repository.NewUserRepository()
	err := userRepo.SaveUser(userModelP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 201, Data: userModelP}, nil
}

func UpdateUserActiveSessionAndToken(userIdP *uint, sessionIdP *uint, tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository()
	err := userRepo.UpdateUser(userIdP, sessionIdP, tokenIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}
