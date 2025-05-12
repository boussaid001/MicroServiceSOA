package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/yourusername/go-microservices-project/services/graphql-service/models"
)

// ReviewRepository defines a repository for review operations
type ReviewRepository struct {
	db *sql.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

// GetAll returns all reviews with optional filtering and pagination
func (r *ReviewRepository) GetAll(params models.ReviewQueryParams) ([]*models.Review, error) {
	// Base query
	query := `
		SELECT id, product_id, user_id, username, rating, comment, created_at
		FROM reviews
		WHERE 1=1
	`
	
	// Add filters if provided
	var args []interface{}
	var conditions []string
	
	argPosition := 1
	
	if params.ProductID != "" {
		conditions = append(conditions, fmt.Sprintf("product_id = $%d", argPosition))
		args = append(args, params.ProductID)
		argPosition++
	}
	
	if params.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argPosition))
		args = append(args, params.UserID)
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
		return nil, fmt.Errorf("failed to query reviews: %w", err)
	}
	defer rows.Close()
	
	// Parse results
	var reviews []*models.Review
	for rows.Next() {
		review := &models.Review{}
		if err := rows.Scan(
			&review.ID,
			&review.ProductID,
			&review.UserID,
			&review.Username,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan review row: %w", err)
		}
		reviews = append(reviews, review)
	}
	
	return reviews, nil
}

// GetByID returns a single review by ID
func (r *ReviewRepository) GetByID(id string) (*models.Review, error) {
	query := `
		SELECT id, product_id, user_id, username, rating, comment, created_at
		FROM reviews
		WHERE id = $1
	`
	
	review := &models.Review{}
	err := r.db.QueryRow(query, id).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.Username,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No review found with this ID
		}
		return nil, fmt.Errorf("failed to query review by ID: %w", err)
	}
	
	return review, nil
}

// GetByProductID returns all reviews for a specific product
func (r *ReviewRepository) GetByProductID(productID string) ([]*models.Review, error) {
	params := models.ReviewQueryParams{
		ProductID: productID,
	}
	return r.GetAll(params)
}

// Create inserts a new review
func (r *ReviewRepository) Create(input models.CreateReviewInput) (*models.Review, error) {
	query := `
		INSERT INTO reviews (product_id, user_id, username, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, product_id, user_id, username, rating, comment, created_at
	`
	
	review := &models.Review{}
	err := r.db.QueryRow(
		query,
		input.ProductID,
		input.UserID,
		input.Username,
		input.Rating,
		input.Comment,
	).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.Username,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}
	
	log.Printf("Created review with ID: %s", review.ID)
	return review, nil
}

// Update updates an existing review
func (r *ReviewRepository) Update(input models.UpdateReviewInput) (*models.Review, error) {
	query := `
		UPDATE reviews
		SET rating = $2, comment = $3
		WHERE id = $1
		RETURNING id, product_id, user_id, username, rating, comment, created_at
	`
	
	review := &models.Review{}
	err := r.db.QueryRow(
		query,
		input.ID,
		input.Rating,
		input.Comment,
	).Scan(
		&review.ID,
		&review.ProductID,
		&review.UserID,
		&review.Username,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to update review: %w", err)
	}
	
	return review, nil
}

// Delete removes a review by ID
func (r *ReviewRepository) Delete(id string) error {
	query := "DELETE FROM reviews WHERE id = $1"
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("review not found")
	}
	
	return nil
} 