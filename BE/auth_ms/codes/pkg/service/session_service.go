package service

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"

	"gorm.io/gorm"
)

type SessionService interface {
	GetSession(sessionIdP *uint, lockForUpdate bool) (any, error)
	GetUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error)
	EndSession(sessionIdP *uint) (any, error)
	StoreSession(userIdP *uint, sessionTokenTraceIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (any, error)
	RefreshSession(sessionModelP *model.Session, tokenIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) error
}

func NewSessionService(newTx *gorm.DB) SessionService {
	if newTx != nil {
		return &baseService{tx: newTx}
	}

	return &baseService{tx: nil}
}

func (s *baseService) GetSession(sessionIdP *uint, lockForUpdate bool) (any, error) {
	sessionRepo := repository.NewSessionRepository(s.tx)

	var sessionModelP *model.Session
	var err error

	if lockForUpdate {
		sessionModelP, err = sessionRepo.FindSessionAndLockForUpdate(sessionIdP)
	} else {
		sessionModelP, err = sessionRepo.FindSession(sessionIdP)
	}

	if err != nil {
		return nil, err
	}

	return sessionModelP, nil
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

	sessionRepo := repository.NewSessionRepository(s.tx)
	err := sessionRepo.CreateSession(sessionModelP)
	if err != nil {
		return nil, err
	}

	return sessionModelP, nil
}

func (s *baseService) RefreshSession(
	sessionModelP *model.Session, tokenIdP *string,
	jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) error {

	*refreshCount++
	sessionModelP.RefreshCount = *refreshCount
	sessionModelP.SessionTokenTraceId = tokenIdP
	sessionModelP.RefreshEndsAt = time.Unix(*refreshExpiresAt, 0)
	sessionModelP.EndsAt = time.Unix(*jwtExpiresAt, 0)

	sessionRepo := repository.NewSessionRepository(s.tx)
	if err := sessionRepo.SaveSession(sessionModelP); err != nil {
		return err
	}

	return nil
}

func (s *baseService) EndSession(sessionIdP *uint) (any, error) {
	sessionRepo := repository.NewSessionRepository(s.tx)
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

func (s *baseService) GetUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error) {
	sessionRepo := repository.NewSessionRepository(s.tx)
	sessionArrayP, err := sessionRepo.FindUserSessions(userIdP, limit, offset)
	if err != nil {
		return nil, err
	}

	return sessionArrayP, nil
}
