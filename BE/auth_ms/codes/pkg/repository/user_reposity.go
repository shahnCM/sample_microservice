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
	UpdateUser(userIdP *uint, updatesP *map[string]any) error
	SaveUser(user *model.User) error
}

func NewUserRepository() UserRepository {
	db := mariadb10.GetMariaDb10()
	return &baseRepository{db: db}
}

func (r *baseRepository) FindUserById(userIdP *uint) (*model.User, error) {
	var user model.User
	// var session model.Session
	if err := r.db.Unscoped().
		Preload("LastSession", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC").Limit(1)
		}).
		Where("id = ?", userIdP).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *baseRepository) FindUser(identifier string, password string) (*model.User, error) {
	var user model.User
	// var session model.Session
	if err := r.db.Unscoped().
		Preload("LastSession", func(db *gorm.DB) *gorm.DB {
			return db.Order("id DESC").Limit(1)
		}).
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

func (r *baseRepository) SaveUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *baseRepository) UpdateUser(userIdP *uint, updatesP *map[string]any) error {
	if err := r.db.Model(&model.User{}).Unscoped().
		Where("id = ?", userIdP).
		Updates(updatesP).
		Error; err != nil {

		return err
	}
	return nil
}
