package product

// GetProductRequest represents a request to get a product by ID
type GetProductRequest struct {
	Id string
}

// ListProductsRequest represents a request to list products
type ListProductsRequest struct {
	Page     int32
	Limit    int32
	Category string
}

// ListProductsResponse represents a response containing products
type ListProductsResponse struct {
	Products []*Product
	Total    int32
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name        string
	Description string
	Price       float32
	Stock       int32
	Category    string
	Images      []string
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Id          string
	Name        string
	Description string
	Price       float32
	Stock       int32
	Category    string
	Images      []string
}

// DeleteProductRequest represents a request to delete a product
type DeleteProductRequest struct {
	Id string
}

// DeleteProductResponse represents a response from deleting a product
type DeleteProductResponse struct {
	Success bool
} 