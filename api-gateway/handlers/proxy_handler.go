package handlers

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewProxyHandler creates a reverse proxy that forwards requests to targetBaseURL
func NewProxyHandler(targetBaseURL string) gin.HandlerFunc {
	target, err := url.Parse(targetBaseURL)
	if err != nil {
		log.Fatalf("Invalid target base URL for proxy: %v", err) 
	}

	// Store target path for endpoint-specific handling
	targetPath := target.Path
	
	// For Hasura, we don't need to modify the target path as it's already included in the target URL
	useTargetPath := strings.Contains(targetBaseURL, "hasura")
	
	if targetPath == "" && !useTargetPath {
		log.Printf("No path provided in target URL, using default path")
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
		
		// For Hasura, use the complete path from the target URL
		if useTargetPath {
			req.URL.Path = targetPath
		}
		
		req.Host = target.Host
		
		log.Printf("Proxying request from %s to -> %s", originalPath, req.URL.String())
	}

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
} 