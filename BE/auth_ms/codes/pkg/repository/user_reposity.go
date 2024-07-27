package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	FindUser(identifier string) (*model.User, error)
	FindUserById(userIdP *uint) (*model.User, error)
	FindUserByIdFast(userIdP *uint) (*model.User, error)
	FindUserByIdAndLockForUpdate(userIdP *uint) (*model.User, error)
	CreateUser(userModelP *model.User) error
	UpdateUser(userModelP *model.User) error
}

func NewUserRepository(tx *gorm.DB) UserRepository {
	if tx != nil {
		return &baseRepository{db: tx}
	}
	db := mariadb10.GetMariaDb10()
	return &baseRepository{db: db}
}

func (r *baseRepository) FindUser(identifier string) (*model.User, error) {
	var user model.User

	if err := r.db.Unscoped().
		// Preload("LastSession", func(db *gorm.DB) *gorm.DB {
		// 	return db.Order("id DESC").Limit(1)
		// }).
		Where("username = ? OR email = ?", identifier, identifier).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *baseRepository) FindUserById(userIdP *uint) (*model.User, error) {
	var user model.User

	if err := r.db.Unscoped().
		// Preload("LastSession", func(db *gorm.DB) *gorm.DB {
		// 	return db.Unscoped().Select("id", "refresh_count").Order("id DESC").Limit(1)
		// }).
		Where("id = ?", userIdP).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *baseRepository) FindUserByIdFast(userIdP *uint) (*model.User, error) {
	var user model.User

	if err := r.db.Unscoped().
		Select("id", "username", "role", "session_token_trace_id").
		Where("id = ?", userIdP).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *baseRepository) FindUserByIdAndLockForUpdate(userIdP *uint) (*model.User, error) {
	var user model.User

	if err := r.db.Unscoped().
		Clauses(clause.Locking{Strength: "UPDATE"}).
		// Preload("LastSession", func(db *gorm.DB) *gorm.DB {
		// 	return db.Unscoped().Select("id", "refresh_count").Order("id DESC").Limit(1)
		// }).
		Where("id = ?", userIdP).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *baseRepository) CreateUser(userModelP *model.User) error {
	return r.db.Create(userModelP).Error
}

func (r *baseRepository) UpdateUser(userModelP *model.User) error {
	return r.db.Unscoped().Save(userModelP).Error
}
