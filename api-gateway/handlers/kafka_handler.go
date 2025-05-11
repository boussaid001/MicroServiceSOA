package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrderHandler handles requests for the Order service
type OrderHandler struct {
	kafkaBrokers string
	producer      sarama.SyncProducer
	orders        map[string]*Order // In-memory store for demo purposes
	mutex         sync.RWMutex
}

// Order represents an order in the system
type Order struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userId"`
	Products   []OrderItem `json:"products"`
	TotalPrice float64    `json:"totalPrice"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(kafkaBrokers string) *OrderHandler {
	// Configure the Kafka producer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Create the Kafka producer
	producer, err := sarama.NewSyncProducer([]string{kafkaBrokers}, config)
	if err != nil {
		// For now, just log the error and continue
		fmt.Printf("Failed to create Kafka producer: %v\n", err)
	}

	return &OrderHandler{
		kafkaBrokers: kafkaBrokers,
		producer:     producer,
		orders:       make(map[string]*Order),
	}
}

// GetOrders returns all orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// Convert map to slice
	orders := make([]Order, 0, len(h.orders))
	for _, order := range h.orders {
		orders = append(orders, *order)
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder returns an order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	id := c.Param("id")
	order, exists := h.orders[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Parse request body
	var req struct {
		UserID   string      `json:"userId"`
		Products []OrderItem `json:"products"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calculate total price
	var totalPrice float64
	for _, item := range req.Products {
		totalPrice += item.Price * float64(item.Quantity)
	}

	// Create new order
	now := time.Now()
	order := &Order{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		Products:   req.Products,
		TotalPrice: totalPrice,
		Status:     "PENDING",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Store order in memory
	h.mutex.Lock()
	h.orders[order.ID] = order
	h.mutex.Unlock()

	// Produce Kafka message
	if h.producer != nil {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize order"})
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: "orders",
			Key:   sarama.StringEncoder(order.ID),
			Value: sarama.ByteEncoder(orderJSON),
		}

		_, _, err = h.producer.SendMessage(msg)
		if err != nil {
			// Log error but continue
			fmt.Printf("Failed to send message to Kafka: %v\n", err)
		}
	}

	c.JSON(http.StatusCreated, order)
}

// UpdateOrder updates an order
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	// Parse request body
	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update order
	h.mutex.Lock()
	defer h.mutex.Unlock()

	order, exists := h.orders[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.Status = req.Status
	order.UpdatedAt = time.Now()

	// Produce Kafka message for order update
	if h.producer != nil {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize order"})
			return
		}

		msg := &sarama.ProducerMessage{
			Topic: "order_updates",
			Key:   sarama.StringEncoder(order.ID),
			Value: sarama.ByteEncoder(orderJSON),
		}

		_, _, err = h.producer.SendMessage(msg)
		if err != nil {
			// Log error but continue
			fmt.Printf("Failed to send message to Kafka: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderStatus gets the status of an order
func (h *OrderHandler) GetOrderStatus(c *gin.Context) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	id := c.Param("id")
	order, exists := h.orders[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     order.ID,
		"status": order.Status,
	})
} 