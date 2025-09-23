package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxRetries  = 3
	retryDelay  = 5 * time.Second
	kafkaBroker = "192.168.1.151:9093"
	kafkaTopic  = "topic-A"
)

func KafkaSend(messageKey, messageValue string) error {
	var lastErr error

	log.Printf("Starting Kafka send. Key: '%s', Value length: %d", messageKey, len(messageValue))

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Waiting %v before retry (attempt %d/%d)...", retryDelay, attempt, maxRetries)
			time.Sleep(retryDelay)
		}

		log.Printf("Attempt %d/%d to connect to Kafka", attempt+1, maxRetries+1)

		w := &kafka.Writer{
			Addr:     kafka.TCP(kafkaBroker),
			Topic:    kafkaTopic,
			Balancer: &kafka.LeastBytes{},
		}

		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(messageKey),
				Value: []byte(messageValue),
			},
		)

		if closeErr := w.Close(); closeErr != nil {
			log.Printf("Warning: error closing writer: %v", closeErr)
		}

		if err == nil {
			log.Printf("Message successfully sent on attempt %d", attempt+1)
			return nil
		}

		lastErr = err
		log.Printf("Attempt %d/%d failed: %v", attempt+1, maxRetries+1, err)

		// Если это последняя попытка - выходим
		if attempt == maxRetries {
			break
		}
	}

	finalErr := fmt.Errorf("failed to send Kafka message after %d attempts. Last error: %w", maxRetries+1, lastErr)
	log.Printf("FATAL: %v", finalErr)
	return finalErr
}
