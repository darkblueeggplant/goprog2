package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	gin.DisableConsoleColor()

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", HealthCheck)
	}

	r.Run(":8080")
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
