package middleware

import (
	"log"
	"net/http"
	"time"
)

// Функция логгер обрабатываемых http-запросов
func Logger(process http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log.Printf("%s:%s recieved", r.Method, r.URL.Path)
		// Обёртка для http.ResponseWriter для получения результата обработки запроса
		stealer := NewResponseStealer(w)
		// Обработка запроса
		process.ServeHTTP(stealer, r)
		level := "[WARN]"
		switch {
		case stealer.StatusCode < 400:
			level = "[INFO]"
		case stealer.StatusCode < 500:
			level = "[ERROR]"
		}
		log.Printf("%s %s %s time - %v status - %d", level, r.Method, r.URL.Path, time.Since(now), stealer.StatusCode)
	})
}
