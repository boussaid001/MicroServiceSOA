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
	// Check if products table exists
	exists, err := p.tableExists("products")
	if err != nil {
		return err
	}

	// Create the products table if it doesn't exist
	if !exists {
		log.Println("Creating products table...")
		err = p.createProductsTable()
		if err != nil {
			return err
		}
		log.Println("Products table created successfully!")
	} else {
		log.Println("Products table already exists, skipping creation")
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

// createProductsTable creates the products table
func (p *PostgresDB) createProductsTable() error {
	query := `
		CREATE TABLE products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
			stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
			category VARCHAR(100),
			images TEXT[],
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX idx_products_category ON products(category);
		CREATE INDEX idx_products_created_at ON products(created_at);
	`
	_, err := p.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}
	return nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
} 