package service

import (
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"
)

func GetSession(sessionIdP *uint) (*response.GenericServiceResponseDto, error) {
	sessionRepo := repository.NewSessionRepository()
	sessionP, err := sessionRepo.FindSession(sessionIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: sessionP}, nil
}

func StoreSession(userIdP *uint, lastTokenIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (*response.GenericServiceResponseDto, error) {
	sessionModelP := &model.Session{
		UserId:        userIdP,
		StartsAt:      time.Now(),
		LastTokenId:   lastTokenIdP,
		EndsAt:        time.Unix(*jwtExpiresAt, 0),
		RefreshEndsAt: time.Unix(*refreshExpiresAt, 0),
		RefreshCount:  0,
	}

	sessionRepo := repository.NewSessionRepository()
	err := sessionRepo.SaveSession(sessionModelP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 201, Data: sessionModelP}, nil
}

func EndSession(sessionIdP *uint) (*response.GenericServiceResponseDto, error) {
	sessionRepo := repository.NewSessionRepository()
	err := sessionRepo.EndSession(sessionIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}
