package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReviewHandler handles requests for the Review service
type ReviewHandler struct {
	baseURL string
}

// NewReviewHandler creates a new ReviewHandler
func NewReviewHandler(baseURL string) *ReviewHandler {
	return &ReviewHandler{
		baseURL: baseURL,
	}
}

// executeGraphQLQuery executes a GraphQL query against the review service
func (h *ReviewHandler) executeGraphQLQuery(query string, variables map[string]interface{}) ([]byte, error) {
	// Create the GraphQL request
	reqBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, err
	}

	// Make the request to the GraphQL endpoint
	resp, err := http.Post(
		fmt.Sprintf("%s/graphql", h.baseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetProductReviews gets reviews for a product
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID := c.Param("productId")

	// GraphQL query to get product reviews
	query := `
		query GetProductReviews($productId: ID!) {
			productReviews(productId: $productId) {
				id
				productId
				userId
				username
				rating
				comment
				createdAt
			}
		}
	`

	variables := map[string]interface{}{
		"productId": productID,
	}

	// Execute the query
	response, err := h.executeGraphQLQuery(query, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pass through the response
	c.Data(http.StatusOK, "application/json", response)
}

// GetReview gets a single review by ID
func (h *ReviewHandler) GetReview(c *gin.Context) {
	reviewID := c.Param("id")

	// GraphQL query to get a single review
	query := `
		query GetReview($id: ID!) {
			review(id: $id) {
				id
				productId
				userId
				username
				rating
				comment
				createdAt
			}
		}
	`

	variables := map[string]interface{}{
		"id": reviewID,
	}

	// Execute the query
	response, err := h.executeGraphQLQuery(query, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pass through the response
	c.Data(http.StatusOK, "application/json", response)
}

// CreateReview creates a new review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	// Parse request body
	var req struct {
		ProductID string  `json:"productId"`
		UserID    string  `json:"userId"`
		Username  string  `json:"username"`
		Rating    float32 `json:"rating"`
		Comment   string  `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// GraphQL mutation to create a review
	mutation := `
		mutation CreateReview($input: CreateReviewInput!) {
			createReview(input: $input) {
				id
				productId
				userId
				username
				rating
				comment
				createdAt
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"productId": req.ProductID,
			"userId":    req.UserID,
			"username":  req.Username,
			"rating":    req.Rating,
			"comment":   req.Comment,
		},
	}

	// Execute the mutation
	response, err := h.executeGraphQLQuery(mutation, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pass through the response
	c.Data(http.StatusCreated, "application/json", response)
}
