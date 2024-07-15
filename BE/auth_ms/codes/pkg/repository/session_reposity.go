package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SessionRepository interface {
	FindSession(identifier *uint) (*model.Session, error)
	SaveSession(session *model.Session) error
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
	if err := r.db.Where("id = ?", identifier).Unscoped().First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) SaveSession(session *model.Session) error {
	return r.db.Create(session).Error
}
