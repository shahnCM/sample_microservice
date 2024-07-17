package model

import (
	"time"
)

// Token represents the tokens table
type Token struct {
	Id               *string   `gorm:"primaryKey;type:char(26);not null"`
	UserId           *uint     `gorm:"not null"`
	SessionId        *uint     `gorm:"not null"`
	JwtToken         *string   `gorm:""`
	RefreshToken     *string   `gorm:""`
	JwtExpiresAt     time.Time `gorm:"not null"`
	RefreshExpiresAt time.Time `gorm:"not null"`
	TokenStatus      string    `gorm:"type:enum('fresh', 'refreshed', 'revoked');not null"`
	Session          Session   `gorm:"foreignKey:SessionId"`
	User             User      `gorm:"foreignKey:UserId"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	// DeletedAt        time.Time
}
