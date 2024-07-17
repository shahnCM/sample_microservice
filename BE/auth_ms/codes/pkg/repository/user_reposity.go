package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserById(userIdP *uint) (*model.User, error)
	FindUser(identifier string, password string) (*model.User, error)
	UpdateUser(userIdP *uint, sessionIdP *uint, tokenIdP *string) error
	SaveUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	db := mariadb10.GetMariaDb10()
	return &userRepository{db: db}
}

func (r *userRepository) FindUserById(userIdP *uint) (*model.User, error) {
	var user model.User
	if err := r.db.Unscoped().
		// Preload("Sessions", "ends_at > ?", time.Now()).
		// Preload("Sessions.Tokens", "jwt_expires_at > ?", time.Now()).
		// Preload("Tokens", "jwt_expires_at > ?", time.Now()).
		Where("id = ?", userIdP).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUser(identifier string, password string) (*model.User, error) {
	var user model.User
	if err := r.db.Unscoped().
		// Preload("Sessions", "ends_at > ?", time.Now()).
		// Preload("Sessions.Tokens", "jwt_expires_at > ?", time.Now()).
		// Preload("Tokens", "jwt_expires_at > ?", time.Now()).
		Where("username = ? OR email = ?", identifier, identifier).
		Where("password = ?", password).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SaveUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(userIdP *uint, sessionIdP *uint, tokenIdP *string) error {
	if err := r.db.Model(&model.User{}).Unscoped().
		Where("id = ?", userIdP).
		Unscoped().
		Update("last_session_id", sessionIdP).
		Update("last_token_id", tokenIdP).Error; err != nil {
		return err
	}
	return nil
}
