package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/1337yeeee/order-service-wb/internal/server"
)

func main() {
	// 1. Инициализация логгера
	logger := log.New(os.Stdout, "ORDER-SERVICE: ", log.LstdFlags|log.Lshortfile)

	// // 2. Инициализация кэша
	// orderCache := cache.New()

	// // 3. Подключение к PostgreSQL
	// dbConn, err := db.Connect()
	// if err != nil {
	// 	logger.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer dbConn.Close()

	// // 4. Восстановление кэша из БД при старте
	// if err := cache.RestoreFromDB(orderCache, dbConn); err != nil {
	// 	logger.Printf("Warning: couldn't restore cache from DB: %v", err)
	// }

	// // 5. Инициализация Kafka Consumer
	// kafkaConsumer, err := kafka.NewConsumer(orderCache, dbConn)
	// if err != nil {
	// 	logger.Fatalf("Failed to create Kafka consumer: %v", err)
	// }
	// go kafkaConsumer.Start() // Запускаем в отдельной горутине

	// 6. Инициализация HTTP сервера
	// httpServer := server.New(orderCache, dbConn)
	httpServer := server.New()
	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 7. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatalf("HTTP server shutdown error: %v", err)
	}

	// kafkaConsumer.Stop()
	logger.Println("Server gracefully stopped")
}
