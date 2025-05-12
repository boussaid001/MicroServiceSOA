package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// NewProxyHandler creates a reverse proxy that forwards requests to targetBaseURL
func NewProxyHandler(targetBaseURL string) gin.HandlerFunc {
	target, err := url.Parse(targetBaseURL)
	if err != nil {
		log.Fatalf("Invalid target base URL for proxy: %v", err) 
	}

	// Check if the target has a path. If not, we'll use the default.
	if target.Path == "" {
		target.Path = "/v1/graphql"
		log.Printf("No path provided in target URL, defaulting to %s", target.Path)
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: target.Scheme,
		Host:   target.Host,
	})

	// Create a custom director to set the target path
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		
		// Store the original path for logging
		originalPath := req.URL.Path
		
		// Always set the path to target.Path (which should be /v1/graphql)
		req.URL.Path = target.Path
		req.Host = target.Host
		
		log.Printf("Proxying request from %s to -> %s", originalPath, req.URL.String())
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
} 