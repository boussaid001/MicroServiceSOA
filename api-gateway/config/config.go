package config

import (
	"os"
)

// Config holds the configuration for the API Gateway
type Config struct {
	Port             string
	RestServiceURL   string
	GrpcServiceURL   string
	GraphqlServiceURL string
	KafkaBrokers     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Port:             getEnv("PORT", "8080"),
		RestServiceURL:   getEnv("REST_SERVICE_URL", "http://localhost:8081"),
		GrpcServiceURL:   getEnv("GRPC_SERVICE_URL", "localhost:8082"),
		GraphqlServiceURL: getEnv("GRAPHQL_SERVICE_URL", "http://localhost:8083"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
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
