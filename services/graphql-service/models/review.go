package models

import (
	"time"
)

// Review represents a product review
type Review struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateReviewInput represents the input data for creating a new review
type CreateReviewInput struct {
	ProductID string  `json:"product_id"`
	UserID    string  `json:"user_id"`
	Username  string  `json:"username"`
	Rating    float64 `json:"rating"`
	Comment   string  `json:"comment"`
}

// UpdateReviewInput represents the input data for updating an existing review
type UpdateReviewInput struct {
	ID       string  `json:"id"`
	Rating   float64 `json:"rating"`
	Comment  string  `json:"comment"`
}

// ReviewQueryParams contains parameters for querying reviews
type ReviewQueryParams struct {
	ProductID string
	UserID    string
	Limit     int
	Offset    int
} 