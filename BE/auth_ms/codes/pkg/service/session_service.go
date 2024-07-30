package service

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/repository"
	"time"

	"gorm.io/gorm"
)

type SessionService interface {
	GetUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error)
	GetSession(sessionIdP *uint, lockForUpdate bool) (*model.Session, error)
	StartSession(userIdP *uint, sessionTokenTraceIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (*model.Session, error)
	RefreshSession(sessionModelP *model.Session, SessionTokenTraceIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) error
	EndSession(sessionIdP *uint) (any, error)
}

func NewSessionService(newTx *gorm.DB) SessionService {
	if newTx != nil {
		return &baseService{tx: newTx}
	}

	return &baseService{tx: nil}
}

func (s *baseService) GetUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error) {
	sessionRepo := repository.NewSessionRepository(s.tx)
	sessionArrayP, err := sessionRepo.FindUserSessions(userIdP, limit, offset)
	if err != nil {
		return nil, err
	}

	return sessionArrayP, nil
}

func (s *baseService) GetSession(sessionIdP *uint, lockForUpdate bool) (*model.Session, error) {
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

func (s *baseService) StartSession(userIdP *uint, sessionTokenTraceIdP *string, jwtExpiresAt *int64, refreshExpiresAt *int64) (*model.Session, error) {
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
	sessionModelP *model.Session, SessionTokenTraceIdP *string,
	jwtExpiresAt *int64, refreshExpiresAt *int64, refreshCount *int) error {

	*refreshCount++
	sessionModelP.RefreshCount = *refreshCount
	sessionModelP.SessionTokenTraceId = SessionTokenTraceIdP
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
	sessionModelP := &model.Session{
		EndsAt:        time.Now(),
		RefreshEndsAt: time.Now(),
		Revoked:       true,
	}
	err := sessionRepo.UpdateSession(sessionIdP, sessionModelP)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
