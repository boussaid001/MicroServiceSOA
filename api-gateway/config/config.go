package config

import (
	"log"
	"os"
)

// Config holds application configuration
type Config struct {
	RestServiceURL    string
	GrpcServiceURL    string
	GraphqlServiceURL string
	HasuraServiceURL  string
	KafkaBrokers      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		RestServiceURL:    getEnv("REST_SERVICE_URL", "http://localhost:8081"),
		GrpcServiceURL:    getEnv("GRPC_SERVICE_URL", "localhost:8082"),
		GraphqlServiceURL: getEnv("GRAPHQL_SERVICE_URL", "http://localhost:8083"),
		HasuraServiceURL:  getEnv("HASURA_SERVICE_URL", "http://localhost:8090/v1/graphql"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "localhost:9092"),
	}

	log.Printf("Loaded configuration: REST=%s, gRPC=%s, GraphQL=%s, Hasura=%s, Kafka=%s",
		cfg.RestServiceURL, cfg.GrpcServiceURL, cfg.GraphqlServiceURL, cfg.HasuraServiceURL, cfg.KafkaBrokers)

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
