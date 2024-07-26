package driver

import (
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLDriver struct{}

func (d MySQLDriver) Connect(dsn string) (*gorm.DB, error) {
	var loggerConfig logger.Interface
	logQueriesToConsole, _ := strconv.ParseBool(os.Getenv("LOG_QUERIES_TO_CONSOLE"))

	if logQueriesToConsole {
		loggerConfig = logger.Default.LogMode(logger.Info)
	} else {
		loggerConfig = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: loggerConfig,
	})

	return db, err
}
