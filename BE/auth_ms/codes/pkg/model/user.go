package model

import (
	"time"
)

// User represents the users table
type User struct {
	Id                  uint     `gorm:"primaryKey;autoIncrement"`
	Username            string   `gorm:"unique;not null"`
	Password            string   `gorm:"not null"`
	Email               string   `gorm:"unique;not null"`
	Role                string   `gorm:"type:enum('standard', 'admin');not null"`
	SessionTokenTraceId *string  `gorm:"unique"`
	LastSessionId       *uint    `gorm:"unique"`
	LastSession         *Session `gorm:"foreignKey:UserId"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time
}
