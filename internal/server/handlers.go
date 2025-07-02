package server

import (
	"log"
	"net/http"
)

func (s *Server) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	const testKey = "test_order"

	// Логирование запроса
	log.Printf("Request received for order: %s", r.URL.Path)

	// Пытаемся получить данные из кэша
	data, ok := s.cache.Get(testKey)
	if !ok {
		http.Error(w, "No data in cache", http.StatusNotFound)
		log.Printf("No data found in cache for key: %s", testKey)
		return
	}

	// Логирование успешного ответа
	log.Printf("Serving data from cache for key: %s", testKey)

	// Устанавливаем заголовки и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
