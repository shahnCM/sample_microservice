package model

import (
	"time"
)

// Session represents the sessions table
type Session struct {
	Id            uint      `gorm:"primaryKey:autoIncrement"`
	UserId        *uint     `gorm:"type:char(26);not null"`
	LastTokenId   *string   `gorm:"type:char(26);not null"`
	RefreshCount  int       `gorm:"not null"`
	StartsAt      time.Time `gorm:"not null"`
	EndsAt        time.Time
	Revoked       bool
	RefreshEndsAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
