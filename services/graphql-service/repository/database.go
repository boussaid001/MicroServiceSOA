package repository

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Database configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// InitDB initializes the database connection
func InitDB(config DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	// Try to connect with retries
	var db *sql.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Failed to open database connection: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Failed to ping database: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Successfully connected
		break
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d retries: %v", maxRetries, err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// EnsureTablesExist creates necessary tables if they don't exist
func EnsureTablesExist(db *sql.DB) error {
	// Read SQL from file
	sqlBytes, err := ioutil.ReadFile("./schema/reviews_table.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %v", err)
	}

	// Execute SQL statements
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %v", err)
	}

	log.Println("Successfully ensured tables exist")
	return nil
} 