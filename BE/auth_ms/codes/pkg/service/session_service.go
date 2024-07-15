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
		return nil, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: sessionP}, nil
}

func StoreSession(userIdP *uint) (*response.GenericServiceResponseDto, error) {
	sessionModelP := &model.Session{
		UserId:       *userIdP,
		StartedAt:    time.Now(),
		RefreshCount: 0,
	}

	sessionRepo := repository.NewSessionRepository()
	err := sessionRepo.SaveSession(sessionModelP)
	if err != nil {
		return nil, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: sessionModelP}, nil
}
