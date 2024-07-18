package service

import (
	"auth_ms/pkg/dto/request"
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/enum"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
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
	user, err := userRepo.FindUser(userP.Username)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: user}, nil
}

func StoreUser(userP *request.UserRegistrationDto) (*response.GenericServiceResponseDto, error) {
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

func UpdateUserActiveToken(userIdP *uint, tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository()
	updatesP := &map[string]any{
		"last_token_id": tokenIdP,
	}
	err := userRepo.UpdateUser(userIdP, updatesP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}

func UpdateUserActiveSessionAndToken(userIdP *uint, sessionIdP *uint, tokenIdP *string) (*response.GenericServiceResponseDto, error) {
	userRepo := repository.NewUserRepository()
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
