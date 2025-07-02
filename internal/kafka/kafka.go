package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/1337yeeee/order-service-wb/internal/models"
	"github.com/1337yeeee/order-service-wb/internal/repository"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader     *kafka.Reader
	repository *repository.OrderRepository
}

func New(brokers []string, topic string, groupID string, repository *repository.OrderRepository) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // Минимум 1 байт для быстрого реагирования
		MaxBytes: 10e6, // Максимум 10MB
	})

	return &Consumer{
		reader:     reader,
		repository: repository,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Получае сообщение из Kafka
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Kafka consumer error: %v", err)
				continue
			}

			// Десериализируем сообщение в модель заказа
			var order models.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}

			// Сохраняем заказ через репозиторий
			if err := c.repository.Save(order); err != nil {
				log.Printf("Failed to save order: %v", err)
				continue
			}
		}
	}
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Printf("Failed to close Kafka reader: %v", err)
	}
}
