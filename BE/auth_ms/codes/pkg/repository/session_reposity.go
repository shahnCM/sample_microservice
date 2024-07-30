package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository interface {
	FindSession(identifier *uint) (*model.Session, error)
	FindSessionAndLockForUpdate(identifier *uint) (*model.Session, error)
	FindUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error)
	CreateSession(sessionModelP *model.Session) error
	SaveSession(sessionModelP *model.Session) error
	UpdateSession(sessionIdP *uint, sessionModelP *model.Session) error
}

func NewSessionRepository(tx *gorm.DB) SessionRepository {
	if tx != nil {
		return &baseRepository{db: tx}
	}
	db := mariadb10.GetMariaDb10()
	return &baseRepository{db: db}
}

func (r *baseRepository) FindSession(identifier *uint) (*model.Session, error) {
	var session model.Session
	if err := r.db.Unscoped().
		Where("id = ?", identifier).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *baseRepository) FindSessionAndLockForUpdate(identifier *uint) (*model.Session, error) {
	var session model.Session
	if err := r.db.Unscoped().
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", identifier).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *baseRepository) FindUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error) {
	var sessionArray []*model.Session

	if err := r.db.Unscoped().
		Model(&model.Session{}).
		Where("user_id = ?", userIdP).
		Order("id").
		Limit(*limit).
		Offset(*offset).
		Scan(&sessionArray).
		Error; err != nil {
		return nil, err
	}
	return &sessionArray, nil
}

func (r *baseRepository) CreateSession(sessionModelP *model.Session) error {
	return r.db.Unscoped().Create(sessionModelP).Error
}

func (r *baseRepository) SaveSession(sessionModelP *model.Session) error {
	return r.db.Unscoped().Save(sessionModelP).Error
}

func (r *baseRepository) UpdateSession(sessionIdP *uint, sessionModelP *model.Session) error {
	if err := r.db.Model(&model.Session{}).Unscoped().
		Where("id = ?", sessionIdP).
		Updates(sessionModelP).
		Error; err != nil {

		return err
	}
	return nil
}
