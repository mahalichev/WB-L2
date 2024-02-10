package handlers

import (
	"encoding/json"
	"net/http"
)

// Функция формирования и отправки http-ответа в формате JSON
func ResponseJSON(w http.ResponseWriter, response any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	switch {
	case statusCode < 400:
		if err := json.NewEncoder(w).Encode(map[string]any{"result": response}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		if err := json.NewEncoder(w).Encode(map[string]any{"error": response}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
