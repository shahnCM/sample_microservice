package service

import "gorm.io/gorm"

type baseService struct {
	tx *gorm.DB
}
