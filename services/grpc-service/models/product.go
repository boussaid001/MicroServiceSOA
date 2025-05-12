package models

import (
	"database/sql"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	Category    string    `json:"category"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductInput represents the input data for creating a new product
type CreateProductInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int32    `json:"stock"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

// UpdateProductInput represents the input data for updating an existing product
type UpdateProductInput struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int32    `json:"stock"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

// ProductQueryParams contains parameters for querying products
type ProductQueryParams struct {
	Category string
	Limit    int
	Offset   int
}

// ScanProduct scans a database row into a Product struct
func ScanProduct(row *sql.Row) (*Product, error) {
	var product Product
	var imagesArray sql.NullString

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.Category,
		&imagesArray,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Handle the images array
	if imagesArray.Valid {
		// Parse the string representation of the array
		// This is a simplified approach - in production you might want to use a more robust solution
		rawImages := imagesArray.String[1 : len(imagesArray.String)-1] // Remove the curly braces
		if rawImages != "" {
			product.Images = parseArrayString(rawImages)
		}
	}

	return &product, nil
}

// ScanProducts scans database rows into a slice of Product
func ScanProducts(rows *sql.Rows) ([]*Product, error) {
	var products []*Product

	for rows.Next() {
		var product Product
		var imagesArray sql.NullString

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Category,
			&imagesArray,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Handle the images array
		if imagesArray.Valid {
			// Parse the string representation of the array
			rawImages := imagesArray.String[1 : len(imagesArray.String)-1] // Remove the curly braces
			if rawImages != "" {
				product.Images = parseArrayString(rawImages)
			}
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// parseArrayString parses a PostgreSQL array string into a slice of strings
func parseArrayString(s string) []string {
	// This is a simplified approach - in production you might want to use a more robust solution
	// that properly handles escaping, quotes, etc.
	result := make([]string, 0)

	// Split by commas
	parts := make([]string, 0)
	inQuote := false
	start := 0

	for i, c := range s {
		if c == '"' {
			inQuote = !inQuote
		} else if c == ',' && !inQuote {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}

	// Add the last part
	if start < len(s) {
		parts = append(parts, s[start:])
	}

	// Remove quotes
	for _, part := range parts {
		part = removeQuotes(part)
		result = append(result, part)
	}

	return result
}

// removeQuotes removes double quotes from the beginning and end of a string if present
func removeQuotes(s string) string {
	s = removeLeadingAndTrailingSpaces(s)
	
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// removeLeadingAndTrailingSpaces removes spaces from the beginning and end of a string
func removeLeadingAndTrailingSpaces(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	s = s[i:]
	
	i = len(s) - 1
	for i >= 0 && (s[i] == ' ' || s[i] == '\t') {
		i--
	}
	if i >= 0 {
		s = s[:i+1]
	}
	
	return s
} 