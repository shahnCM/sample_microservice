package service

import (
	"auth_ms/pkg/dto/response"
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"

	"gorm.io/gorm"
)

func GetSession(tx *gorm.DB, sessionIdP *uint) (*response.GenericServiceResponseDto, error) {
	sessionRepo := repository.NewSessionRepository(tx)
	sessionP, err := sessionRepo.FindSession(sessionIdP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 404, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 200, Data: sessionP}, nil
}

func StoreSession(tx *gorm.DB, userIdP *uint, lastTokenIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (*response.GenericServiceResponseDto, error) {
	sessionModelP := &model.Session{
		UserId:        userIdP,
		StartsAt:      time.Now(),
		LastTokenId:   lastTokenIdP,
		EndsAt:        time.Unix(*jwtExpiresAt, 0),
		RefreshEndsAt: time.Unix(*refreshExpiresAt, 0),
		RefreshCount:  0,
	}

	sessionRepo := repository.NewSessionRepository(tx)
	err := sessionRepo.SaveSession(sessionModelP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 201, Data: sessionModelP}, nil
}

func RefreshSession(tx *gorm.DB, sessionIdP *uint, userIdP *uint, tokenIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) (*response.GenericServiceResponseDto, error) {
	*refreshCount++
	sessionRepo := repository.NewSessionRepository(tx)
	updatesP := &map[string]any{
		"last_token_id":   tokenIdP,
		"ends_at":         time.Unix(*jwtExpiresAt, 0),
		"refresh_ends_at": time.Unix(*refreshExpiresAt, 0),
		"refresh_count":   refreshCount,
	}
	err := sessionRepo.UpdateSession(sessionIdP, updatesP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}

func EndSession(tx *gorm.DB, sessionIdP *uint) (*response.GenericServiceResponseDto, error) {
	sessionRepo := repository.NewSessionRepository(tx)
	updatesP := &map[string]interface{}{
		"ends_at":         time.Now(),
		"refresh_ends_at": time.Now(),
		"revoked":         true,
	}
	err := sessionRepo.UpdateSession(sessionIdP, updatesP)
	if err != nil {
		return &response.GenericServiceResponseDto{StatusCode: 422, Data: nil}, err
	}

	return &response.GenericServiceResponseDto{StatusCode: 204, Data: nil}, nil
}
