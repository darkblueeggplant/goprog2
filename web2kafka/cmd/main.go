package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	Time      string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	ClientIP  string `json:"client_ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Query     string `json:"query,omitempty"`
	Protocol  string `json:"protocol"`
	Status    int    `json:"status_code"`
	LatencyMs int64  `json:"latency_ms"`
	Latency   string `json:"latency"`
	UserAgent string `json:"user_agent,omitempty"`
	Referer   string `json:"referer,omitempty"`
	BytesIn   int64  `json:"bytes_in,omitempty"`
	BytesOut  int    `json:"bytes_out,omitempty"`
	Error     string `json:"error,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func main() {
	r := gin.Default()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		entry := LogEntry{
			Time:      param.TimeStamp.Format(time.RFC3339),
			Level:     "info",
			Service:   "gin-api",
			ClientIP:  param.ClientIP,
			Method:    param.Method,
			Path:      param.Path,
			Query:     param.Request.URL.RawQuery,
			Protocol:  param.Request.Proto,
			Status:    param.StatusCode,
			LatencyMs: param.Latency.Milliseconds(),
			Latency:   param.Latency.String(),
			UserAgent: param.Request.UserAgent(),
			Referer:   param.Request.Referer(),
			BytesIn:   param.Request.ContentLength,
			BytesOut:  param.BodySize,
			Error:     param.ErrorMessage,
		}

		jsonBytes, _ := json.Marshal(entry)
		return string(jsonBytes) + "\n"
	}))

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
