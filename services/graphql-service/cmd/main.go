package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func main() {
	// Get port from environment or use default
	port := getEnv("PORT", "8083")

	// Create a simple schema for reviews
	schema, err := createSchema()
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}

	// Create a GraphQL HTTP handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Set up Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// GraphQL endpoint
	router.POST("/graphql", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	// GraphiQL UI
	router.GET("/graphql", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("GraphQL Review Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Create a simple GraphQL schema
func createSchema() (graphql.Schema, error) {
	// Define review type
	reviewType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Review",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
				},
				"productId": &graphql.Field{
					Type: graphql.String,
				},
				"userId": &graphql.Field{
					Type: graphql.String,
				},
				"username": &graphql.Field{
					Type: graphql.String,
				},
				"rating": &graphql.Field{
					Type: graphql.Float,
				},
				"comment": &graphql.Field{
					Type: graphql.String,
				},
				"createdAt": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	// Root query type
	rootQuery := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"review": &graphql.Field{
					Type: reviewType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.ID,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// Mock implementation for now
						id, ok := p.Args["id"].(string)
						if ok {
							// Return mock data
							return map[string]interface{}{
								"id":        id,
								"productId": "product-123",
								"userId":    "user-456",
								"username":  "johndoe",
								"rating":    4.5,
								"comment":   "Great product!",
								"createdAt": "2023-01-01T12:00:00Z",
							}, nil
						}
						return nil, nil
					},
				},
				"productReviews": &graphql.Field{
					Type: graphql.NewList(reviewType),
					Args: graphql.FieldConfigArgument{
						"productId": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						// Mock implementation for now
						productID, ok := p.Args["productId"].(string)
						if ok {
							// Return mock data
							return []map[string]interface{}{
								{
									"id":        "review-1",
									"productId": productID,
									"userId":    "user-456",
									"username":  "johndoe",
									"rating":    4.5,
									"comment":   "Great product!",
									"createdAt": "2023-01-01T12:00:00Z",
								},
								{
									"id":        "review-2",
									"productId": productID,
									"userId":    "user-789",
									"username":  "janedoe",
									"rating":    5.0,
									"comment":   "Excellent service!",
									"createdAt": "2023-01-02T10:30:00Z",
								},
							}, nil
						}
						return nil, nil
					},
				},
			},
		},
	)

	// Schema configuration
	schemaConfig := graphql.SchemaConfig{
		Query: rootQuery,
	}

	// Create and return the schema
	return graphql.NewSchema(schemaConfig)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
