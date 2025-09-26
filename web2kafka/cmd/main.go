package main

import (
	"log/slog"
	"net/http"
	"web2kafka/producer"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 1. Отправляем сообщение о получении запроса со всеми деталями
		go func() {
			message := "Request: " + c.Request.Method + " " + c.Request.URL.Path +
				" from " + clientIP +
				" | Agent: " + c.Request.UserAgent() +
				" | Query: " + c.Request.URL.RawQuery

			err := producer.KafkaSend("request-received", message)
			if err != nil {
				slog.Error("Failed to send to Kafka", "error", err)
			} else {
				slog.Info("Message sent to Kafka successfully")
			}
		}()

		// 2. Отвечаем клиенту
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello!",
			"status":  "success",
		})

		// 3. Отправляем сообщение об отправке ответа
		go func() {
			message := "Response: 200 OK to " + clientIP +
				" | Method: " + c.Request.Method +
				" | Path: " + c.Request.URL.Path

			err := producer.KafkaSend("response-sent", message)
			if err != nil {
				slog.Error("Failed to send to Kafka", "error", err)
			} else {
				slog.Info("Second message sent to Kafka successfully")
			}
		}()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	slog.Info("Starting server on :8080")
	r.Run(":8080")
}
