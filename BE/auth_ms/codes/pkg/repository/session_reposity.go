package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SessionRepository interface {
	FindSession(identifier *uint) (*model.Session, error)
	SaveSession(session *model.Session) error
	EndSession(sessionIdP *uint) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository() SessionRepository {
	db := mariadb10.GetMariaDb10()
	return &sessionRepository{db: db}
}

func (r *sessionRepository) FindSession(identifier *uint) (*model.Session, error) {
	var session model.Session
	if err := r.db.Unscoped().
		// Preload("Tokens").
		Where("id = ?", identifier).
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) SaveSession(session *model.Session) error {
	return r.db.Unscoped().Create(session).Error
}

func (r *sessionRepository) EndSession(sessionIdP *uint) error {
	if err := r.db.Model(&model.Session{}).Unscoped().
		Where("id = ?", sessionIdP).
		Where("ends_at > ?", time.Now()).
		Update("ends_at", time.Now()).
		Update("refresh_ends_at", time.Now()).
		Error; err != nil {

		return err
	}
	return nil
}
