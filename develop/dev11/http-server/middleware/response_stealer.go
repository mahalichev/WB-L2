package middleware

import "net/http"

// Обертка для http.ResponseWriter для получения результата обработки http-запроса
type ResponseStealer struct {
	http.ResponseWriter
	StatusCode int
}

func (stealer *ResponseStealer) WriteHeader(statusCode int) {
	stealer.StatusCode = statusCode
	stealer.ResponseWriter.WriteHeader(statusCode)
}

func NewResponseStealer(w http.ResponseWriter) *ResponseStealer {
	return &ResponseStealer{w, http.StatusOK}
}
