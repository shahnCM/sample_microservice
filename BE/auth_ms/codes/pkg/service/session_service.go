package service

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"

	"gorm.io/gorm"
)

var tx *gorm.DB = nil

func SetTx(newtx *gorm.DB) {
	tx = newtx
}

func UnsetTx() {
	tx = nil
}

type SessionService interface {
	GetSession(sessionIdP *uint) (any, error)
	StoreSession(userIdP *uint, sessionTokenTraceIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (any, error)
	RefreshSession(sessionIdP *uint, userIdP *uint, tokenIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) (any, error)
	EndSession(sessionIdP *uint) (any, error)
}

func NewSessionService(newTx *gorm.DB) SessionService {
	if tx != nil {
		return &baseService{tx: newTx}
	}

	return &baseService{tx: nil}
}

func (s *baseService) GetSession(sessionIdP *uint) (any, error) {
	sessionRepo := repository.NewSessionRepository(s.tx)
	sessionP, err := sessionRepo.FindSession(sessionIdP)
	if err != nil {
		return nil, err
	}

	return sessionP, nil
}

func (s *baseService) StoreSession(userIdP *uint, sessionTokenTraceIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (any, error) {
	sessionModelP := &model.Session{
		UserId:              userIdP,
		StartsAt:            time.Now(),
		SessionTokenTraceId: sessionTokenTraceIdP,
		EndsAt:              time.Unix(*jwtExpiresAt, 0),
		RefreshEndsAt:       time.Unix(*refreshExpiresAt, 0),
		RefreshCount:        0,
	}

	sessionRepo := repository.NewSessionRepository(tx)
	err := sessionRepo.SaveSession(sessionModelP)
	if err != nil {
		return nil, err
	}

	return sessionModelP, nil
}

func (s *baseService) RefreshSession(sessionIdP *uint, userIdP *uint, tokenIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) (any, error) {
	*refreshCount++
	sessionRepo := repository.NewSessionRepository(tx)
	updatesP := &map[string]any{
		"session_token_trace_id": tokenIdP,
		"ends_at":                time.Unix(*jwtExpiresAt, 0),
		"refresh_ends_at":        time.Unix(*refreshExpiresAt, 0),
		"refresh_count":          refreshCount,
	}
	err := sessionRepo.UpdateSession(sessionIdP, updatesP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *baseService) EndSession(sessionIdP *uint) (any, error) {
	sessionRepo := repository.NewSessionRepository(tx)
	updatesP := &map[string]interface{}{
		"ends_at":         time.Now(),
		"refresh_ends_at": time.Now(),
		"revoked":         true,
	}
	err := sessionRepo.UpdateSession(sessionIdP, updatesP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
