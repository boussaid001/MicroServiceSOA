package product

import (
	"context"
)

// This is a placeholder file to make this a valid Go package
// Normally, this would contain the compiled protobuf code
// For a real implementation, you would run protoc to generate
// the Go code from the .proto file

// Product represents a product in the catalog
type Product struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Stock       int32    `json:"stock"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

// ProductServiceClient is the client API for ProductService service
type ProductServiceClient interface {
	// GetProduct retrieves a product by ID
	GetProduct(ctx context.Context, req *GetProductRequest) (*Product, error)
	
	// ListProducts retrieves a list of products with optional filtering
	ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error)
	
	// CreateProduct creates a new product
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
	
	// UpdateProduct updates an existing product
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*Product, error)
	
	// DeleteProduct deletes a product by ID
	DeleteProduct(ctx context.Context, req *DeleteProductRequest) (*DeleteProductResponse, error)
}

// ProductService is the interface for the product service
type ProductService interface {
	GetProduct(id string) (*Product, error)
	ListProducts(page, limit int32, category string) ([]*Product, int32, error)
	CreateProduct(product *Product) (*Product, error)
	UpdateProduct(product *Product) (*Product, error)
	DeleteProduct(id string) (bool, error)
} 