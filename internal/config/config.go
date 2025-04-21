package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	AppName              string
	AppEnv               string
	DatabaseUser         string
	DatabasePassword     string
	DatabaseHost         string
	DatabaseName         string
	DatabasePort         string
	DatabaseMaxOpenConns int
	DatabaseMaxIdleConns int
	Port                 string
}

const (
	defaultDatabaseMaxOpenConns = 10
	defaultDatabaseMaxIdleConns = 10
)

func Load() *Config {
	environment := getEnv("APP_ENV", "development")

	err := godotenv.Load()
	if err != nil && environment == "development" {
		zap.L().Info("No .env file found, using environment variables")
	}

	config := &Config{
		AppName:              getEnv("APP_NAME", "fiap-sa-order-service"),
		AppEnv:               environment,
		DatabaseUser:         fetchEnv("DB_USER"),
		DatabasePassword:     getEnv("DB_PASSWORD", ""),
		DatabaseHost:         fetchEnv("DB_HOST"),
		DatabaseName:         fetchEnv("DB_NAME"),
		DatabasePort:         fetchEnv("DB_PORT"),
		DatabaseMaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", defaultDatabaseMaxOpenConns),
		DatabaseMaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", defaultDatabaseMaxIdleConns),
		Port:                 getEnv("PORT", "8082"),
	}

	return config
}

func (cfg Config) IsDevelopment() bool {
	return cfg.AppEnv == "development"
}

func (cfg Config) IsProduction() bool {
	return cfg.AppEnv == "production"
}
