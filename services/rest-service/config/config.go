package config

import (
	"os"
)

// Config holds all configuration for the service
type Config struct {
	Port     string
	Database DatabaseConfig
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Port: getEnv("PORT", "8081"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "userdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
