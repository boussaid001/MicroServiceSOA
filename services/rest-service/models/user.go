package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't expose password in JSON responses
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email" binding:"omitempty,email"`
	Password  string `json:"password" binding:"omitempty,min=6"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// UserResponse represents the user data returned in API responses
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
