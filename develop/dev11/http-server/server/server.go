package server

import (
	"context"
	"dev11/http-server/api"
	"dev11/http-server/handlers"
	"dev11/http-server/middleware"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Создание маршрутизатора http-запросов
func newMux(manager *api.EventManager) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/create_event", handlers.CreateEvent(manager))
	mux.Handle("/update_event", handlers.UpdateEvent(manager))
	mux.Handle("/delete_event", handlers.DeleteEvent(manager))
	mux.Handle("/events_for_day", handlers.GetEvents(manager, "day"))
	mux.Handle("/events_for_week", handlers.GetEvents(manager, "week"))
	mux.Handle("/events_for_month", handlers.GetEvents(manager, "month"))
	// Оборачивание маршрутизатора в функцию-логгер (декоратор)
	return middleware.Logger(mux)
}

// Инициализация и запуск сервера
func StartServer(address, port string) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", address, port),
		Handler: newMux(api.NewEventManager()),
	}
	// Запуск сервера
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	UntilInterrupt()

	// Корректное выключение сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Ожидание сигнала о завершении работы программы
func UntilInterrupt() {
	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)
	<-c
}
