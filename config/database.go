package config

import (
	"errors"
	"os"
)

var errorMissingDatabaseConfiguration = errors.New("missing database configuration")

type DatabaseConfig struct {
	Url      string
	Username string
	Password string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	url := os.Getenv("DB_URL")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	if url == "" || username == "" || password == "" {
		return nil, errorMissingDatabaseConfiguration
	}

	return &DatabaseConfig{
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}
