package product

import (
	"context"

	"google.golang.org/grpc"
)

// ProductServiceServer is the server API for ProductService service.
type ProductServiceServer interface {
	GetProduct(context.Context, *GetProductRequest) (*Product, error)
	ListProducts(context.Context, *ListProductsRequest) (*ListProductsResponse, error)
	CreateProduct(context.Context, *CreateProductRequest) (*Product, error)
	UpdateProduct(context.Context, *UpdateProductRequest) (*Product, error)
	DeleteProduct(context.Context, *DeleteProductRequest) (*DeleteProductResponse, error)
}

// RegisterProductService registers the ProductService with the gRPC server
func RegisterProductService(s *grpc.Server, srv ProductServiceServer) {
	// Register the service with the gRPC server
	// In a real implementation, this would use the generated code from protoc
} 