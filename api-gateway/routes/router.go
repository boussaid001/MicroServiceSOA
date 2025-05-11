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
	reviewHandler := handlers.NewReviewHandler(cfg.GraphqlServiceURL)

	// API group
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

		// Review routes (GraphQL)
		reviews := api.Group("/reviews")
		{
			reviews.GET("/product/:productId", reviewHandler.GetProductReviews)
			reviews.GET("/:id", reviewHandler.GetReview)
			reviews.POST("/", reviewHandler.CreateReview)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
