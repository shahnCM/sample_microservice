package model

import (
	"time"
)

// Session represents the sessions table
type Session struct {
	Id            uint      `gorm:"primaryKey:autoIncrement"`
	StartsAt      time.Time `gorm:"not null"`
	EndsAt        time.Time `gorm:""`
	RefreshEndsAt time.Time `gorm:""`
	LastTokenId   *string   `gorm:"type:char(26);not null"`
	RefreshCount  int       `gorm:"not null"`
	UserId        *uint     `gorm:"type:char(26);not null"`
	Tokens        []Token   `gorm:"foreignKey:SessionId"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	// DeletedAt    time.Time
}
