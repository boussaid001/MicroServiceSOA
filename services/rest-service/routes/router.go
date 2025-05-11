package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-microservices-project/services/rest-service/controllers"
	"github.com/yourusername/go-microservices-project/services/rest-service/repository"
)

// SetupRouter creates and configures a Gin router
func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// Create repositories
	userRepo := repository.NewUserRepository(db)

	// Create controllers
	userController := controllers.NewUserController(userRepo)

	// Set up CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// User routes
	users := router.Group("/users")
	{
		users.GET("/", userController.GetUsers)
		users.GET("/:id", userController.GetUser)
		users.POST("/", userController.CreateUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	return router
}
