package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	pb "github.com/yourusername/go-microservices-project/proto/product"
)

// ProductHandler handles requests for the Product service
type ProductHandler struct {
	serviceURL string
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(serviceURL string) *ProductHandler {
	return &ProductHandler{
		serviceURL: serviceURL,
	}
}

// connect returns a mock gRPC client for the product service
func (h *ProductHandler) connect() (pb.ProductServiceClient, error) {
	// For now, just create a mock client
	client := pb.NewProductServiceClient(nil)
	return client, nil
}

// GetProducts returns a list of products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Connect to gRPC service
	client, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	// Make gRPC request
	resp, err := client.ListProducts(context.Background(), &pb.ListProductsRequest{
		Page:     int32(page),
		Limit:    int32(limit),
		Category: category,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProduct returns a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Connect to gRPC service
	client, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get product ID from URL
	id := c.Param("id")

	// Make gRPC request
	resp, err := client.GetProduct(context.Background(), &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Connect to gRPC service
	client, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse request body
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float32  `json:"price"`
		Stock       int32    `json:"stock"`
		Category    string   `json:"category"`
		Images      []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make gRPC request
	resp, err := client.CreateProduct(context.Background(), &pb.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      req.Images,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Connect to gRPC service
	client, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get product ID from URL
	id := c.Param("id")

	// Parse request body
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float32  `json:"price"`
		Stock       int32    `json:"stock"`
		Category    string   `json:"category"`
		Images      []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make gRPC request
	resp, err := client.UpdateProduct(context.Background(), &pb.UpdateProductRequest{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		Images:      req.Images,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteProduct deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Connect to gRPC service
	client, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get product ID from URL
	id := c.Param("id")

	// Make gRPC request
	resp, err := client.DeleteProduct(context.Background(), &pb.DeleteProductRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
