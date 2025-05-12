package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-microservices-project/api-gateway/config"
	"github.com/yourusername/go-microservices-project/api-gateway/handlers"
)

// SetupRoutes sets up all the routes for the API Gateway
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Create handlers
	userHandler := handlers.NewUserHandler(cfg.RestServiceURL)
	productHandler := handlers.NewProductHandler(cfg.GrpcServiceURL)
	orderHandler := handlers.NewOrderHandler(cfg.KafkaBrokers)
	// reviewHandler := handlers.NewReviewHandler(cfg.GraphqlServiceURL) // Keep for now, might be used for other review-related REST endpoints if any

	// Create proxies for the GraphQL services
	customGraphqlProxyHandler := handlers.NewProxyHandler(cfg.GraphqlServiceURL)
	hasuraProxyHandler := handlers.NewProxyHandler(cfg.HasuraServiceURL)

	// API group for versioning or other common path prefix (optional here for /graphql)
	// Example: api := router.Group("/api")

	// GraphQL routes
	router.POST("/graphql", customGraphqlProxyHandler)
	router.GET("/graphql", customGraphqlProxyHandler) // For GraphiQL access
	
	// Hasura GraphQL routes
	router.POST("/hasura", hasuraProxyHandler)
	router.GET("/hasura", hasuraProxyHandler) // For Hasura console/GraphiQL access

	// Existing API group for RESTful services
	api := router.Group("/api")
	{
		// User routes (REST)
		users := api.Group("/users")
		{
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("/", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Product routes (gRPC)
		products := api.Group("/products")
		{
			products.GET("/", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.POST("/", productHandler.CreateProduct)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}

		// Order routes (Kafka)
		orders := api.Group("/orders")
		{
			orders.GET("/", orderHandler.GetOrders)
			orders.GET("/:id", orderHandler.GetOrder)
			orders.POST("/", orderHandler.CreateOrder)
			orders.PUT("/:id", orderHandler.UpdateOrder)
			orders.GET("/status/:id", orderHandler.GetOrderStatus)
		}

		// Commenting out old /api/reviews if they are fully replaced by /graphql
		// reviews := api.Group("/reviews")
		// {
		// 	reviews.GET("/product/:productId", reviewHandler.GetProductReviews)
		// 	reviews.GET("/:id", reviewHandler.GetReview)
		// 	reviews.POST("/", reviewHandler.CreateReview)
		// }
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
