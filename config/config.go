package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string
	DBName     string
	DBHost     string
	DBUser     string
	DBPass     string
	DBPort     string
	JWTSecret  string
}

func validateConfig(cfg *Config) {
	if cfg.ServerPort == "" {
		log.Fatal("SERVER_PORT is required")
	}
	if cfg.DBHost == "" {
		log.Fatal("DB_HOST is required")
	}
	if cfg.DBUser == "" {
		log.Fatal("DB_USER is required")
	}
	if cfg.DBName == "" {
		log.Fatal("DB_NAME is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
}

func LoadConfig() *Config {
	cfg := &Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPass:     os.Getenv("DB_PASS"),
		DBPort:     os.Getenv("DB_PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	validateConfig(cfg)
	return cfg
}
