package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TokenRepository interface {
	FindToken(identifier *string) (*model.Token, error)
	SaveToken(token *model.Token) error
	UpdateTokenStatus(identifier *string, tokenStatus string) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository() TokenRepository {
	db := mariadb10.GetMariaDb10()
	return &tokenRepository{db: db}
}

func (r *tokenRepository) FindToken(identifier *string) (*model.Token, error) {
	var token model.Token
	if err := r.db.Unscoped().
		// Preload("Session").
		Where("id = ?", identifier).
		First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepository) SaveToken(token *model.Token) error {
	return r.db.Unscoped().Create(token).Error
}

func (r *tokenRepository) UpdateTokenStatus(identifier *string, tokenStatus string) error {
	if err := r.db.Unscoped().Model(&model.Token{}).
		Where("id = ?", identifier).
		Update("token_status", tokenStatus).Error; err != nil {
		return err
	}
	return nil
}
