package model

import (
	"time"

	"gorm.io/gorm"
)

// Token represents the tokens table
type Token struct {
	gorm.Model
	Id               *string `gorm:"primaryKey;type:char(26);not null"`
	UserId           *uint   `gorm:"not null"`
	SessionId        *uint   `gorm:"not null"`
	JwtToken         *string `gorm:""`
	RefreshToken     *string `gorm:""`
	JwtExpiresAt     *int64  `gorm:"not null"`
	RefreshExpiresAt *int64  `gorm:"not null"`
	TokenStatus      string  `gorm:"type:enum('fresh', 'refreshed', 'revoked');not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
