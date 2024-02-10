package models

import (
	"errors"
	"time"
)

var ErrNonPositiveID error = errors.New("event ID must be positive")
var ErrNonPositiveUserID error = errors.New("user ID must be positive")
var ErrEmptyTitle error = errors.New("title cannot be empty")

type Event struct {
	Date        time.Time
	Title       string
	Description string
	ID          int
	UserID      int
}

func NewEvent(id, userId int, date time.Time, title, description string) Event {
	return Event{
		ID:          id,
		UserID:      userId,
		Date:        date,
		Title:       title,
		Description: description,
	}
}

// Валидация события
func ValidateEvent(event Event) error {
	// Если ID события не положительное число - ошибка
	if event.ID <= 0 {
		return ErrNonPositiveID
	}
	// Если ID пользователя не положительное число - ошибка
	if event.UserID <= 0 {
		return ErrNonPositiveUserID
	}
	// Если название события пустое - ошибка
	if event.Title == "" {
		return ErrEmptyTitle
	}
	return nil
}
