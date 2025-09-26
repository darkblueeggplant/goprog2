package producer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxRetries  = 3
	retryDelay  = 5 * time.Second
	kafkaBroker = "192.168.1.151:9093"
	kafkaTopic  = "logs"
)

// KafkaSend отправляет сообщение в Kafka с использованием slog для логирования
func KafkaSend(messageKey, messageValue string) error {
	var lastErr error

	slog.Info("Starting Kafka send",
		"key", messageKey,
		"value_length", len(messageValue),
		"topic", kafkaTopic,
		"broker", kafkaBroker,
	)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			slog.Warn("Waiting before retry",
				"delay", retryDelay,
				"attempt", attempt,
				"max_attempts", maxRetries,
			)
			time.Sleep(retryDelay)
		}

		slog.Info("Attempting Kafka connection",
			"attempt", attempt+1,
			"total_attempts", maxRetries+1,
		)

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
			slog.Warn("Error closing Kafka writer", "error", closeErr)
		}

		if err == nil {
			slog.Info("Message successfully sent",
				"attempt", attempt+1,
				"key", messageKey,
			)
			return nil
		}

		lastErr = err
		slog.Error("Kafka send attempt failed",
			"attempt", attempt+1,
			"total_attempts", maxRetries+1,
			"error", err,
		)

		if attempt == maxRetries {
			break
		}
	}

	finalErr := fmt.Errorf("failed to send Kafka message after %d attempts. Last error: %w", maxRetries+1, lastErr)
	slog.Error("All Kafka send attempts failed",
		"error", finalErr,
		"max_attempts", maxRetries+1,
		"key", messageKey,
	)
	return finalErr
}
