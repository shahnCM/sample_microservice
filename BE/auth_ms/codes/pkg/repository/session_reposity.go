package repository

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"

	"gorm.io/gorm"
)

type SessionRepository interface {
	FindSession(identifier *uint) (*model.Session, error)
	FindUserSessions(userIdP *string, limit, offset *int) (*[]*model.Session, error)
	SaveSession(session *model.Session) error
	UpdateSession(sessionIdP *uint, updates *map[string]any) error
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

func (r *baseRepository) SaveSession(session *model.Session) error {
	return r.db.Unscoped().Create(session).Error
}

func (r *baseRepository) UpdateSession(sessionIdP *uint, updatesP *map[string]any) error {
	if err := r.db.Model(&model.Session{}).Unscoped().
		Where("id = ?", sessionIdP).
		Updates(updatesP).
		Error; err != nil {

		return err
	}
	return nil
}
