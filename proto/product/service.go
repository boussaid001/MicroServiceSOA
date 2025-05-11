package product

import (
	"context"
	"time"
)

// MockProductService provides a mock implementation of ProductServiceClient
type MockProductService struct{}

// NewProductServiceClient creates a new ProductServiceClient
func NewProductServiceClient(cc interface{}) ProductServiceClient {
	return &MockProductService{}
}

// GetProduct implements ProductServiceClient
func (s *MockProductService) GetProduct(ctx context.Context, req *GetProductRequest) (*Product, error) {
	// Return a mock product
	return &Product{
		Id:          req.Id,
		Name:        "Mock Product",
		Description: "This is a mock product for testing",
		Price:       99.99,
		Stock:       100,
		Category:    "mock",
		Images:      []string{"https://example.com/image.jpg"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}, nil
}

// ListProducts implements ProductServiceClient
func (s *MockProductService) ListProducts(ctx context.Context, req *ListProductsRequest) (*ListProductsResponse, error) {
	// Return mock products
	products := []*Product{
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

	return &ListProductsResponse{
		Products: products,
		Total:    int32(len(products)),
	}, nil
}

// CreateProduct implements ProductServiceClient
func (s *MockProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error) {
	// Return a mock product with data from the request
	return &Product{
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

// UpdateProduct implements ProductServiceClient
func (s *MockProductService) UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*Product, error) {
	// Return a mock product with data from the request
	return &Product{
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

// DeleteProduct implements ProductServiceClient
func (s *MockProductService) DeleteProduct(ctx context.Context, req *DeleteProductRequest) (*DeleteProductResponse, error) {
	// Return a mock response
	return &DeleteProductResponse{
		Success: true,
	}, nil
} 