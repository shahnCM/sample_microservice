package config

import (
	"fmt"
	"os"
)

type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
	Version   string
}

type MariaDb10Config struct {
	Host                  string
	Port                  string
	Name                  string
	User                  string
	Password              string
	Driver                string
	DSN                   string
	MaxOpenConnection     string
	MaxIdleConnection     string
	MaxConnectionLifetime string
}

func GetElasticsearchConfig() *ElasticsearchConfig {
	return &ElasticsearchConfig{
		Addresses: []string{os.Getenv("ES_URL")},
		Username:  os.Getenv("ES_USERNAME"),
		Password:  os.Getenv("ES_PASSWORD"),
		Version:   os.Getenv("ES_VERSION"),
	}
}

func GetMariaDb10Config() *MariaDb10Config {
	host := os.Getenv("MARIA_DB_HOST")
	port := os.Getenv("MARIA_DB_PORT")
	name := os.Getenv("MARIA_DB_NAME")
	user := os.Getenv("MARIA_DB_USER")
	password := os.Getenv("MARIA_DB_PASSWORD")
	driver := os.Getenv("MARIA_DB_DRIVER")
	maxOpenConn := os.Getenv("MARIA_DB_MAX_OPEN_CONNECTION")
	maxIdleConn := os.Getenv("MARIA_DB_MAX_IDLE_CONNECTION")
	maxConnLifetime := os.Getenv("MARIA_DB_MAX_CONNECTION_LIFETIME")

	return &MariaDb10Config{
		Driver:                driver,
		Host:                  host,
		Port:                  port,
		Name:                  name,
		User:                  user,
		Password:              password,
		MaxOpenConnection:     maxOpenConn,
		MaxIdleConnection:     maxIdleConn,
		MaxConnectionLifetime: maxConnLifetime,
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, name),
	}
}
