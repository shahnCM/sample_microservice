package model

import (
	"time"

	"gorm.io/gorm"
)

// Session represents the sessions table
type Session struct {
	gorm.Model
	Id           uint      `gorm:"primaryKey:autoIncrement"`
	StartedAt    time.Time `gorm:"not null"`
	ExpiredAt    time.Time `gorm:""`
	RefreshCount int       `gorm:"not null"`
	UserId       uint      `gorm:"type:char(26);not null"`
	Tokens       []Token   `gorm:"foreignKey:SessionId"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
