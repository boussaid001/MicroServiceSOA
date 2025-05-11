package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/go-microservices-project/services/rest-service/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetAll retrieves all users from the database
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
	SELECT id, username, email, password, first_name, last_name, created_at, updated_at
	FROM users
	ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
	SELECT id, username, email, password, first_name, last_name, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Create adds a new user to the database
func (r *UserRepository) Create(user models.CreateUserRequest) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	query := `
	INSERT INTO users (username, email, password, first_name, last_name, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, username, email, password, first_name, last_name, created_at, updated_at
	`

	now := time.Now()
	var newUser models.User
	err = r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		string(hashedPassword),
		user.FirstName,
		user.LastName,
		now,
		now,
	).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.Password,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// Update updates an existing user
func (r *UserRepository) Update(id int, user models.UpdateUserRequest) (*models.User, error) {
	// Get the existing user first
	existingUser, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update the user fields if provided
	if user.Username != "" {
		existingUser.Username = user.Username
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		existingUser.Password = string(hashedPassword)
	}
	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}

	// Update in database
	query := `
	UPDATE users
	SET username = $1, email = $2, password = $3, first_name = $4, last_name = $5, updated_at = $6
	WHERE id = $7
	RETURNING id, username, email, password, first_name, last_name, created_at, updated_at
	`

	now := time.Now()
	var updatedUser models.User
	err = r.db.QueryRow(
		query,
		existingUser.Username,
		existingUser.Email,
		existingUser.Password,
		existingUser.FirstName,
		existingUser.LastName,
		now,
		id,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.Password,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
