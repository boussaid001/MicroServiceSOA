package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/boussaid001/go-microservices-project/proto"
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

// connect establishes a connection to the gRPC service
func (h *ProductHandler) connect() (pb.ProductServiceClient, context.Context, *grpc.ClientConn, error) {
	// Set timeout for gRPC client connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	
	// Connect to the gRPC server
	conn, err := grpc.Dial(h.serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, err
	}
	
	// Create a gRPC client
	client := pb.NewProductServiceClient(conn)
	return client, ctx, conn, nil
}

// GetProducts returns a list of products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Connect to gRPC service
	client, ctx, conn, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	// Make gRPC request
	resp, err := client.ListProducts(ctx, &pb.ListProductsRequest{
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
	client, ctx, conn, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// Get product ID from URL
	id := c.Param("id")

	// Make gRPC request
	resp, err := client.GetProduct(ctx, &pb.GetProductRequest{
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
	client, ctx, conn, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

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
	resp, err := client.CreateProduct(ctx, &pb.CreateProductRequest{
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
	client, ctx, conn, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

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
	resp, err := client.UpdateProduct(ctx, &pb.UpdateProductRequest{
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
	client, ctx, conn, err := h.connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// Get product ID from URL
	id := c.Param("id")

	// Make gRPC request
	resp, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
