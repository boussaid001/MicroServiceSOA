package resolvers

import (
	"time"

	"github.com/yourusername/go-microservices-project/services/graphql-service/models"
	"github.com/yourusername/go-microservices-project/services/graphql-service/repository"
)

// Resolver is the root resolver
type Resolver struct {
	reviewRepo *repository.ReviewRepository
}

// NewResolver creates a new root resolver
func NewResolver(reviewRepo *repository.ReviewRepository) *Resolver {
	return &Resolver{
		reviewRepo: reviewRepo,
	}
}

// Review represents a review in the GraphQL schema
type Review struct {
	model *models.Review
}

// ID returns the review ID
func (r *Review) ID() string {
	return r.model.ID
}

// ProductID returns the product ID
func (r *Review) ProductID() string {
	return r.model.ProductID
}

// UserID returns the user ID
func (r *Review) UserID() string {
	return r.model.UserID
}

// Username returns the username
func (r *Review) Username() string {
	return r.model.Username
}

// Rating returns the review rating
func (r *Review) Rating() float64 {
	return r.model.Rating
}

// Comment returns the review comment
func (r *Review) Comment() *string {
	if r.model.Comment == "" {
		return nil
	}
	return &r.model.Comment
}

// CreatedAt returns the creation time as ISO-8601 string
func (r *Review) CreatedAt() string {
	return r.model.CreatedAt.Format(time.RFC3339)
}

// CreateReviewInput represents the input for creating a review
type CreateReviewInput struct {
	ProductID string  `json:"productId"`
	UserID    string  `json:"userId"`
	Username  string  `json:"username"`
	Rating    float64 `json:"rating"`
	Comment   *string `json:"comment"`
}

// UpdateReviewInput represents the input for updating a review
type UpdateReviewInput struct {
	ID      string  `json:"id"`
	Rating  float64 `json:"rating"`
	Comment *string `json:"comment"`
}

// ReviewsArgs represents arguments for the reviews query
type ReviewsArgs struct {
	ProductID *string
	UserID    *string
	Limit     *int32
	Offset    *int32
}

// ReviewArgs represents arguments for the review query
type ReviewArgs struct {
	ID string
}

// Reviews resolves the reviews query
func (r *Resolver) Reviews(args ReviewsArgs) ([]*Review, error) {
	params := models.ReviewQueryParams{}

	if args.ProductID != nil {
		params.ProductID = *args.ProductID
	}

	if args.UserID != nil {
		params.UserID = *args.UserID
	}

	if args.Limit != nil {
		params.Limit = int(*args.Limit)
	}

	if args.Offset != nil {
		params.Offset = int(*args.Offset)
	}

	reviews, err := r.reviewRepo.GetAll(params)
	if err != nil {
		return nil, err
	}

	var result []*Review
	for _, review := range reviews {
		result = append(result, &Review{model: review})
	}

	return result, nil
}

// Review resolves the review query
func (r *Resolver) Review(args ReviewArgs) (*Review, error) {
	review, err := r.reviewRepo.GetByID(args.ID)
	if err != nil {
		return nil, err
	}

	if review == nil {
		return nil, nil
	}

	return &Review{model: review}, nil
}

// CreateReview resolves the createReview mutation
func (r *Resolver) CreateReview(args struct{ Input CreateReviewInput }) (*Review, error) {
	input := models.CreateReviewInput{
		ProductID: args.Input.ProductID,
		UserID:    args.Input.UserID,
		Username:  args.Input.Username,
		Rating:    args.Input.Rating,
	}

	if args.Input.Comment != nil {
		input.Comment = *args.Input.Comment
	}

	review, err := r.reviewRepo.Create(input)
	if err != nil {
		return nil, err
	}

	return &Review{model: review}, nil
}

// UpdateReview resolves the updateReview mutation
func (r *Resolver) UpdateReview(args struct{ Input UpdateReviewInput }) (*Review, error) {
	input := models.UpdateReviewInput{
		ID:     args.Input.ID,
		Rating: args.Input.Rating,
	}

	if args.Input.Comment != nil {
		input.Comment = *args.Input.Comment
	}

	review, err := r.reviewRepo.Update(input)
	if err != nil {
		return nil, err
	}

	return &Review{model: review}, nil
}

// DeleteReview resolves the deleteReview mutation
func (r *Resolver) DeleteReview(args struct{ ID string }) (bool, error) {
	err := r.reviewRepo.Delete(args.ID)
	if err != nil {
		return false, err
	}

	return true, nil
} 