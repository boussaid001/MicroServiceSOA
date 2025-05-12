package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/yourusername/go-microservices-project/services/graphql-service/repository"
	"github.com/yourusername/go-microservices-project/services/graphql-service/resolvers"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting GraphQL Review Service...")

	// Load configuration from environment variables
	dbConfig := repository.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "reviewdb"),
	}

	// Initialize database connection
	db, err := repository.InitDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Ensure tables exist
	if err := repository.EnsureTablesExist(db); err != nil {
		log.Fatalf("Failed to ensure tables exist: %v", err)
	}

	// Create repositories
	reviewRepo := repository.NewReviewRepository(db)

	// Read schema
	schemaFile := "./schema/schema.graphql"
	schemaBytes, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}
	schemaString := string(schemaBytes)

	// Create resolver
	resolver := resolvers.NewResolver(reviewRepo)

	// Create schema
	schema := graphql.MustParseSchema(schemaString, resolver)

	// Set up HTTP handler
	http.Handle("/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GraphQL Review Service\n"))
		w.Write([]byte("Use /graphql endpoint for queries and mutations"))
	})))

	http.Handle("/graphql", corsMiddleware(&relay.Handler{Schema: schema}))

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8083",
		Handler: nil, // Use default ServeMux
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8083")
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

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass to next handler
		next.ServeHTTP(w, r)
	})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 