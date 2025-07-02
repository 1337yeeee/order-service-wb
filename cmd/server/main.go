package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1337yeeee/order-service-wb/internal/cache"
	"github.com/1337yeeee/order-service-wb/internal/kafka"
	"github.com/1337yeeee/order-service-wb/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "ORDER-SERVICE: ", log.LstdFlags|log.Lshortfile)

	orderCache := cache.New()

	kafkaConsumer, err := kafka.New(
		[]string{os.Getenv("KAFKA_BROKERS")},
		os.Getenv("KAFKA_TOPIC"),
		os.Getenv("KAFKA_GROUP_ID"),
		orderCache,
	)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Запуск consumer в фоне
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaConsumer.Start(ctx)

	// Инициализация HTTP сервера
	httpServer := server.New(orderCache)
	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	httpServer.Shutdown(shutdownCtx)
	kafkaConsumer.Close()
	logger.Println("Server stopped")
}
