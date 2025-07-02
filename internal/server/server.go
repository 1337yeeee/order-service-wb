package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/1337yeeee/order-service-wb/internal/cache"
)

type Server struct {
	httpServer *http.Server
	cache      *cache.Cache
}

func New(cache *cache.Cache) *Server {
	return &Server{
		cache: cache,
	}
}

func (s *Server) Start() error {
	// Настройка маршрутов
	router := http.NewServeMux()
	router.HandleFunc("/order/", s.getOrderHandler) // Обработчик для заказов

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
