package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUser(identifier string, password string) (*model.User, error)
	SaveUser(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	db := mariadb10.GetMariaDb10()
	return &userRepository{db: db}
}

func (r *userRepository) FindUser(identifier string, password string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ? OR email = ?", identifier, identifier).Where("password = ?", password).Unscoped().First(&user).Error; err != nil {
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
