package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/yourusername/go-microservices-project/proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ProductService implements the gRPC product service
type ProductService struct {
	pb.ProductServiceClient
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	return &pb.Product{
		Id:          req.Id,
		Name:        "Sample Product",
		Description: "This is a sample product from the gRPC service",
		Price:       99.99,
		Stock:       100,
		Category:    "electronics",
		Images:      []string{"https://example.com/image.jpg"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// ListProducts retrieves a list of products
func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products := []*pb.Product{
		{
			Id:          "1",
			Name:        "Product 1",
			Description: "Description of product 1",
			Price:       99.99,
			Stock:       100,
			Category:    req.Category,
			Images:      []string{"https://example.com/image1.jpg"},
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		},
		{
			Id:          "2",
			Name:        "Product 2",
			Description: "Description of product 2",
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

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	return &pb.Product{
		Id:          "new-id",
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

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
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

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

func main() {
	// Get port from environment or use default
	port := getEnv("PORT", "8082")

	// Create a TCP listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register ProductService
	pb.RegisterProductService(grpcServer, &ProductService{})

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	// Start server
	log.Printf("gRPC Product Service starting on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
