package repository

import "gorm.io/gorm"

type baseRepository struct {
	db *gorm.DB
}
