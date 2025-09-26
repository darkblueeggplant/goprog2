package main

import (
	"log"
	"net/http"
	"web2kafka/producer"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 1. Отправляем сообщение о получении запроса
		go func() {
			message := "Hello request received from " + clientIP
			err := producer.KafkaSend("request-received", message)
			if err != nil {
				log.Printf("Failed to send to Kafka: %v", err)
			} else {
				log.Printf("Message sent to Kafka successfully")
			}
		}()

		// 2. Всегда отвечаем "Hello!" без имени
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello!", // Всегда просто "Hello!"
			"status":  "success",
		})

		// 3. Отправляем сообщение об отправке ответа
		go func() {
			message := "Response sent to client: " + clientIP
			err := producer.KafkaSend("response-sent", message)
			if err != nil {
				log.Printf("Failed to send to Kafka: %v", err)
			} else {
				log.Printf("Second message sent to Kafka successfully")
			}
		}()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
