package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (s *Server) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем order_uid из URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		log.Printf("Invalid request path: %s", r.URL.Path)
		return
	}

	orderUID := pathParts[2]
	log.Printf("Request received for order: %s", orderUID)

	// Получаем заказ из репозитория
	order, err := s.repository.Get(orderUID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		log.Printf("Failed to get order: %v", err)
		return
	}

	// Сериализуем Order в JSON
	jsonData, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Failed to marshal order: %v", err)
		return
	}

	// Устанавливаем заголовки и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}

	log.Printf("Successfully served order: %s", orderUID)
}
