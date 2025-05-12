package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/boussaid001/go-microservices-project/services/grpc-service/models"
)

// ProductRepository defines a repository for product operations
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// GetAll returns all products with optional filtering and pagination
func (r *ProductRepository) GetAll(params models.ProductQueryParams) ([]*models.Product, error) {
	// Base query
	query := `
		SELECT id, name, description, price, stock, category, images, created_at, updated_at
		FROM products
		WHERE 1=1
	`
	
	// Add filters if provided
	var args []interface{}
	var conditions []string
	
	argPosition := 1
	
	if params.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argPosition))
		args = append(args, params.Category)
		argPosition++
	}
	
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	
	// Add ordering
	query += " ORDER BY created_at DESC"
	
	// Add pagination
	if params.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPosition)
		args = append(args, params.Limit)
		argPosition++
		
		if params.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argPosition)
			args = append(args, params.Offset)
		}
	}
	
	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()
	
	// Parse results
	products, err := models.ScanProducts(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to scan products: %w", err)
	}
	
	return products, nil
}

// GetByID returns a single product by ID
func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, stock, category, images, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	
	row := r.db.QueryRow(query, id)
	product, err := models.ScanProduct(row)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No product found with this ID
		}
		return nil, fmt.Errorf("failed to query product by ID: %w", err)
	}
	
	return product, nil
}

// GetByCategory returns all products for a specific category
func (r *ProductRepository) GetByCategory(category string) ([]*models.Product, error) {
	params := models.ProductQueryParams{
		Category: category,
	}
	return r.GetAll(params)
}

// Create inserts a new product
func (r *ProductRepository) Create(input models.CreateProductInput) (*models.Product, error) {
	query := `
		INSERT INTO products (name, description, price, stock, category, images)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, price, stock, category, images, created_at, updated_at
	`
	
	row := r.db.QueryRow(
		query,
		input.Name,
		input.Description,
		input.Price,
		input.Stock,
		input.Category,
		pq.Array(input.Images),
	)
	
	product, err := models.ScanProduct(row)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	
	log.Printf("Created product with ID: %s", product.ID)
	return product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(input models.UpdateProductInput) (*models.Product, error) {
	query := `
		UPDATE products
		SET name = $2, 
			description = $3, 
			price = $4, 
			stock = $5, 
			category = $6, 
			images = $7,
			updated_at = $8
		WHERE id = $1
		RETURNING id, name, description, price, stock, category, images, created_at, updated_at
	`
	
	row := r.db.QueryRow(
		query,
		input.ID,
		input.Name,
		input.Description,
		input.Price,
		input.Stock,
		input.Category,
		pq.Array(input.Images),
		time.Now(),
	)
	
	product, err := models.ScanProduct(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	
	return product, nil
}

// Delete removes a product by ID
func (r *ProductRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	
	return nil
} 