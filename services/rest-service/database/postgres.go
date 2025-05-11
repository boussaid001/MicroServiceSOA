package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/yourusername/go-microservices-project/services/rest-service/config"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection with a ping
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create the users table if it doesn't exist
	if err := initDB(db); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

// initDB creates the necessary tables if they don't exist
func initDB(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(createTableSQL)
	return err
}
