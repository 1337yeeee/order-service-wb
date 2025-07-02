package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1337yeeee/order-service-wb/internal/cache"
	"github.com/1337yeeee/order-service-wb/internal/db"
	"github.com/1337yeeee/order-service-wb/internal/kafka"
	"github.com/1337yeeee/order-service-wb/internal/repository"
	"github.com/1337yeeee/order-service-wb/internal/server"
)

func main() {
	// Инициализация логгера
	logger := log.New(os.Stdout, "ORDER-SERVICE: ", log.LstdFlags|log.Lshortfile)

	// Инициализация кэша
	orderCache := cache.New()

	// Инициализация соединения с базой данных
	dbConn, err := db.Connect(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	))
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Инициализация репозитория
	orderRepository := repository.New(orderCache, dbConn)

	// Инициализация Kafka consumer
	kafkaConsumer, err := kafka.New(
		[]string{os.Getenv("KAFKA_BROKERS")},
		os.Getenv("KAFKA_TOPIC"),
		os.Getenv("KAFKA_GROUP_ID"),
		orderRepository,
	)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Запускаем Kafka-потребителя в фоне
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go kafkaConsumer.Start(ctx)

	// Инициализация HTTP сервера
	httpServer := server.New(orderRepository)
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
