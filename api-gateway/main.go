package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-microservices-project/api-gateway/config"
	"github.com/yourusername/go-microservices-project/api-gateway/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set up routes
	routes.SetupRoutes(router, cfg)

	// Serve static files for frontend - place this after the API routes to avoid conflicts
	router.Static("/static", "./frontend/static")
	
	// Redirect root to dashboard instead of serving index.html directly
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	})
	
	// Serve all the static HTML files
	router.StaticFile("/dashboard", "./frontend/index.html")
	router.StaticFile("/users.html", "./frontend/users.html")
	router.StaticFile("/products.html", "./frontend/products.html")
	router.StaticFile("/reviews.html", "./frontend/reviews.html")
	router.StaticFile("/orders.html", "./frontend/orders.html")
	router.StaticFile("/compare-graphql.html", "./frontend/compare-graphql.html")

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting API Gateway server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
