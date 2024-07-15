package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table
type User struct {
	gorm.Model
	Id        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Role      string    `gorm:"type:enum('standard', 'admin');not null"`
	Sessions  []Session `gorm:"foreignKey:UserId"`
	Tokens    []Token   `gorm:"foreignKey:UserId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
