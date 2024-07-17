package service

import (
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserById(tx *gorm.DB, userIdP *uint) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository(tx)
	user, err := userRepo.FindUserById(userIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: user}, nil
}

func GetUser(tx *gorm.DB, userP *request.UserLoginDto) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository(tx)
	user, err := userRepo.FindUser(userP.Username, userP.Password)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: user}, nil
}

func StoreUser(tx *gorm.DB, userP *request.UserRegistrationDto) (*response.GenericServiceResponseDto, error) {
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

	userRepo := repository.NewUserRepository(tx)
	err := userRepo.SaveUser(userModelP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 201, Data: userModelP}, nil
}

func UpdateUserActiveToken(tx *gorm.DB, userIdP *uint, tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository(tx)
	updatesP := &map[string]any{
		"last_token_id": tokenIdP,
	}
	err := userRepo.UpdateUser(userIdP, updatesP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}

func UpdateUserActiveSessionAndToken(tx *gorm.DB, userIdP *uint, sessionIdP *uint, tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository(tx)
	updatesP := &map[string]any{
		"last_session_id": sessionIdP,
		"last_token_id":   tokenIdP,
	}
	err := userRepo.UpdateUser(userIdP, updatesP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}
