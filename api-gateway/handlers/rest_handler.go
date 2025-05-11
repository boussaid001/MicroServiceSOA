package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles requests for the User service
type UserHandler struct {
	baseURL string
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(baseURL string) *UserHandler {
	return &UserHandler{
		baseURL: baseURL,
	}
}

// GetUsers returns a list of users
func (h *UserHandler) GetUsers(c *gin.Context) {
	resp, err := http.Get(fmt.Sprintf("%s/users", h.baseURL))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

// GetUser returns a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", h.baseURL, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/users", h.baseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/users/%s", h.baseURL, id),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/users/%s", h.baseURL, id),
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}
