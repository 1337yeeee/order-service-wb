package kafka

import (
	"context"
	"log"

	"github.com/1337yeeee/order-service-wb/internal/cache"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	cache  *cache.Cache
}

func New(brokers []string, topic string, groupID string, cache *cache.Cache) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // Минимум 1 байт для быстрого реагирования
		MaxBytes: 10e6, // Максимум 10MB
	})

	return &Consumer{
		reader: reader,
		cache:  cache,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Kafka consumer error: %v", err)
				continue
			}

			const testKey = "test_order"
			c.cache.Set(testKey, msg.Value)
			log.Printf(
				"Message stored in cache (topic=%s partition=%d offset=%d)",
				msg.Topic, msg.Partition, msg.Offset,
			)
		}
	}
}

func (c *Consumer) Close() {
	if err := c.reader.Close(); err != nil {
		log.Printf("Failed to close Kafka reader: %v", err)
	}
}
