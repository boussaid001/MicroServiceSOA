package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/boussaid001/go-microservices-project/proto"
)

// Product represents a product in the database
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	Category    string    `json:"category"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductService implements the gRPC product service
type ProductService struct {
	pb.UnimplementedProductServiceServer
	db *sql.DB
}

// GetProduct handles the GetProduct gRPC request
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	// For now, return a mock product
	return &pb.Product{
		Id:          req.Id,
		Name:        "Mock Product",
		Description: "This is a mock product",
		Price:       99.99,
		Stock:       100,
		Category:    "Mock",
		Images:      []string{"https://example.com/image.jpg"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// ListProducts handles the ListProducts gRPC request
func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// For now, return mock products
	products := []*pb.Product{
		{
			Id:          "1",
			Name:        "Mock Product 1",
			Description: "This is the first mock product",
			Price:       99.99,
			Stock:       100,
			Category:    req.Category,
			Images:      []string{"https://example.com/image1.jpg"},
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			Id:          "2",
			Name:        "Mock Product 2",
			Description: "This is the second mock product",
			Price:       199.99,
			Stock:       50,
			Category:    req.Category,
			Images:      []string{"https://example.com/image2.jpg"},
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
	}

	return &pb.ListProductsResponse{
		Products: products,
		Total:    int32(len(products)),
	}, nil
}

// CreateProduct handles the CreateProduct gRPC request
func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	// Generate a unique ID (using a random UUID format)
	id := fmt.Sprintf("%x-%x-%x-%x-%x", 
		time.Now().UnixNano()&0xffffffff, 
		time.Now().UnixNano()>>32&0xffff, 
		time.Now().UnixNano()>>48&0x0fff|0x4000, 
		time.Now().UnixNano()&0x3fff|0x8000, 
		time.Now().UnixNano()&0xffffffffffff)

	// Return a product with the generated ID
	return &pb.Product{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      req.Images,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// UpdateProduct handles the UpdateProduct gRPC request
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	// For now, return a mock updated product
	return &pb.Product{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      req.Images,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// DeleteProduct handles the DeleteProduct gRPC request
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	// For now, return success
	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

func main() {
	log.Println("Starting gRPC Product Service...")

	// Get database connection parameters from environment variables
	dbHost := getEnv("DB_HOST", "postgres-product")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "productdb")
	sslMode := getEnv("DB_SSLMODE", "disable")

	// Create database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping database, retrying in 2 seconds... (attempt %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not ping database after 5 attempts: %v", err)
	}

	log.Println("Successfully connected to database")

	// Ensure products table exists
	if err := ensureProductsTableExists(db); err != nil {
		log.Fatalf("Failed to ensure tables exist: %v", err)
	}

	// Create gRPC server
	port := getEnv("PORT", "8082")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register service and reflection
	productService := &ProductService{db: db}
	pb.RegisterProductServiceServer(grpcServer, productService)
	reflection.Register(grpcServer)

	// Start server in a goroutine
	go func() {
		log.Printf("gRPC Product Service starting on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server exiting")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// tableExists checks if a given table exists in the database
func tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		);
	`
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if table exists: %w", err)
	}
	return exists, nil
}

// createProductsTable creates the products table
func createProductsTable(db *sql.DB) error {
	query := `
		CREATE TABLE products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
			stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
			category VARCHAR(100),
			images TEXT[],
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX idx_products_category ON products(category);
		CREATE INDEX idx_products_created_at ON products(created_at);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}
	return nil
}

// ensureProductsTableExists checks if products table exists and creates it if it doesn't
func ensureProductsTableExists(db *sql.DB) error {
	// Check if products table exists
	exists, err := tableExists(db, "products")
	if err != nil {
		return err
	}

	// Create the products table if it doesn't exist
	if !exists {
		log.Println("Creating products table...")
		err = createProductsTable(db)
		if err != nil {
			return err
		}
		log.Println("Products table created successfully!")
	} else {
		log.Println("Products table already exists, skipping creation")
	}

	return nil
}
