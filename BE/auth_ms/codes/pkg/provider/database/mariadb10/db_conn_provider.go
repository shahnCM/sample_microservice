package mariadb10

import (
	"auth_ms/pkg/config"
	"auth_ms/pkg/driver"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

var txMode bool = false
var db *gorm.DB
var dbOg *gorm.DB
var once sync.Once

func ConnectToMariaDb10() error {

	var connectError error

	mariaDb10river := driver.MySQLDriver{}
	mariaDb10Config := config.GetMariaDb10Config()
	dsn := mariaDb10Config.DSN

	maxOpenConn, err := strconv.Atoi(mariaDb10Config.MaxOpenConnection)
	if err != nil {
		return err
		// panic(fiber.NewError(fiber.StatusInternalServerError, "failed to parse MaxOpenConnection: "+err.Error()))
	}

	maxIdleConn, err := strconv.Atoi(mariaDb10Config.MaxIdleConnection)
	if err != nil {
		return err
		// panic(fiber.NewError(fiber.StatusInternalServerError, "failed to parse MaxIdleConnection: "+err.Error()))
	}

	maxConnLifetime, err := time.ParseDuration(mariaDb10Config.MaxConnectionLifetime)
	if err != nil {
		return err
		// panic(fiber.NewError(fiber.StatusInternalServerError, "failed to parse MaxConnectionLifetime: "+err.Error()))
	}

	once.Do(func() {
		connectError = nil

		db, err = mariaDb10river.Connect(dsn)
		dbOg = db

		if err != nil {
			connectError = err
			// panic(fiber.NewError(fiber.StatusInternalServerError, "failed to connect database"))
		}

		sqlDb, err := db.DB()
		if err != nil {
			connectError = err
			// panic(fiber.NewError(fiber.StatusInternalServerError, "failed to get MariaDb10 instance"))
		}

		sqlDb.SetMaxOpenConns(maxOpenConn)
		sqlDb.SetMaxIdleConns(maxIdleConn)
		sqlDb.SetConnMaxLifetime(maxConnLifetime)
	})

	return connectError
}

func GetMariaDb10() *gorm.DB {
	return db
}

func TransactionBegin() *gorm.DB {
	txMode = true
	db = db.Begin()
	return db
}

func TransactionCommit() *gorm.DB {
	if txMode {
		db.Commit()
		txMode = false
		db = dbOg
	}
	return db
}

func TransactionRollback() *gorm.DB {
	if txMode {
		db.Rollback()
		txMode = false
		db = dbOg
	}
	return db
}
