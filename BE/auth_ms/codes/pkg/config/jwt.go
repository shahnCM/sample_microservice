package config

import (
	"os"
)

type JwtConfig struct {
	JwtSecret        string
	JwtExpiresIn     string
	RefreshExpiresIn string
}

func GetJwtConfig() *JwtConfig {
	return &JwtConfig{
		JwtSecret:        os.Getenv("JWT_SECRET"),
		JwtExpiresIn:     os.Getenv("JWT_EXPIRES_IN"),
		RefreshExpiresIn: os.Getenv("REFRESH_EXPIRES_IN"),
	}
}
