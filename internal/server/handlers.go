package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Пока заглушка - будем возвращать тестовые данные
	orderID := r.URL.Path[len("/order/"):] // Извлекаем ID из URL

	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Заменить на реальные данные из кэша/БД
	testOrder := map[string]interface{}{
		"id":      orderID,
		"status":  "processed",
		"message": "This is a test response. Real implementation will come later!",
		"warning": "Kafka and database integration is not implemented yet",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testOrder)
}
