package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/1337yeeee/order-service-wb/internal/repository"
)

type Server struct {
	httpServer *http.Server
	repository *repository.OrderRepository
}

func New(repository *repository.OrderRepository) *Server {
	return &Server{
		repository: repository,
	}
}

func (s *Server) Start() error {
	// Настройка маршрутов
	router := http.NewServeMux()

	router.Handle("/", http.FileServer(http.Dir("./web")))

	// Обработчик API-запросов
	router.HandleFunc("/order/", s.getOrderHandler)

	// Конфигурация HTTP сервера
	s.httpServer = &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
