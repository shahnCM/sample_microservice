package driver

import "gorm.io/gorm"

type DriverInterface interface {
	Connect(dsn string) (*gorm.DB, error)
}
