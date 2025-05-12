package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// PostgresDB wraps a sql.DB instance
type PostgresDB struct {
	*sql.DB
}

// NewPostgresDB creates a new PostgresDB instance
func NewPostgresDB(dsn string) (*PostgresDB, error) {
	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping database, retrying in 2 seconds... (attempt %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("could not ping database after 5 attempts: %w", err)
	}

	// Configure the database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	pgDB := &PostgresDB{DB: db}

	return pgDB, nil
}

// EnsureTablesExist checks if necessary tables exist and creates them if they don't
func (p *PostgresDB) EnsureTablesExist() error {
	// Check if reviews table exists
	exists, err := p.tableExists("reviews")
	if err != nil {
		return err
	}

	// Create the reviews table if it doesn't exist
	if !exists {
		log.Println("Creating reviews table...")
		err = p.createReviewsTable()
		if err != nil {
			return err
		}
		log.Println("Reviews table created successfully!")
	} else {
		log.Println("Reviews table already exists, skipping creation")
	}

	return nil
}

// tableExists checks if a given table exists in the database
func (p *PostgresDB) tableExists(tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		);
	`
	err := p.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if table exists: %w", err)
	}
	return exists, nil
}

// createReviewsTable creates the reviews table
func (p *PostgresDB) createReviewsTable() error {
	query := `
		CREATE TABLE reviews (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id VARCHAR(255) NOT NULL,
			user_id VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			rating DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 5),
			comment TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX idx_reviews_product_id ON reviews(product_id);
		CREATE INDEX idx_reviews_user_id ON reviews(user_id);
		CREATE INDEX idx_reviews_created_at ON reviews(created_at);
	`
	_, err := p.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create reviews table: %w", err)
	}
	return nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
} 