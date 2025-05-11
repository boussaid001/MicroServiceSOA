package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-microservices-project/api-gateway/config"
	"github.com/yourusername/go-microservices-project/api-gateway/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	// Enable CORS
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

	// Set up routes
	routes.SetupRoutes(router, cfg)

	// Serve static files for frontend - place this after the API routes to avoid conflicts
	router.Static("/static", "./frontend/static")
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/users.html", "./frontend/users.html")
	router.StaticFile("/products.html", "./frontend/products.html")
	router.StaticFile("/reviews.html", "./frontend/reviews.html")
	router.StaticFile("/orders.html", "./frontend/orders.html")

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("API Gateway starting on port %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
