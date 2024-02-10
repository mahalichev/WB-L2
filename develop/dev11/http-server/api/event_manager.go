package api

import (
	"errors"
	"slices"
	"time"

	"dev11/http-server/models"
)

var ErrEventNotFound error = errors.New("event with given ID not found")
var ErrUserNotFound error = errors.New("user with given ID not found")

type EventManager struct {
	events      map[int]map[int]models.Event
	nextEventID int
}

// Функция создания события в календаре
func (manager *EventManager) CreateEvent(userId int, date time.Time, title, description string) (models.Event, error) {
	event := models.NewEvent(manager.nextEventID, userId, date, title, description)
	// Если данные не валидны - ошибка
	if err := models.ValidateEvent(event); err != nil {
		return models.Event{}, err
	}

	// Добавление события в календарь
	if _, ok := manager.events[event.UserID]; !ok {
		manager.events[event.UserID] = make(map[int]models.Event)
	}
	manager.events[event.UserID][event.ID] = event
	manager.nextEventID++
	return event, nil
}

// Функция обновления события в календаре
func (manager *EventManager) UpdateEvent(eventId, userId int, date time.Time, title, description string) (models.Event, error) {
	event := models.NewEvent(eventId, userId, date, title, description)
	// Если данные не валидны - ошибка
	if err := models.ValidateEvent(event); err != nil {
		return models.Event{}, err
	}
	// Если пользователь не найден - ошибка
	if _, ok := manager.events[userId]; !ok {
		return models.Event{}, ErrUserNotFound
	}
	// Если событие не найдено - ошибка
	if _, ok := manager.events[userId][eventId]; !ok {
		return models.Event{}, ErrEventNotFound
	}
	// Обновление события
	manager.events[userId][eventId] = event
	return event, nil
}

// Функция удаления события из календаря
func (manager *EventManager) DeleteEvent(eventId, userId int) error {
	// Если пользователь не найден - ошибка
	if _, ok := manager.events[userId]; !ok {
		return ErrUserNotFound
	}
	// Если событие не найдено - ошибка
	if _, ok := manager.events[userId][eventId]; !ok {
		return ErrEventNotFound
	}
	// Удаление события
	delete(manager.events[userId], eventId)
	return nil
}

// Функция получения списка событий пользователя из календаря в интервале since <= date < to
func (manager *EventManager) GetEvents(userId int, since, to time.Time) ([]models.Event, error) {
	// Если пользователь не найден - ошибка
	if _, ok := manager.events[userId]; !ok {
		return nil, ErrUserNotFound
	}
	result := make([]models.Event, 0)
	// Получение событий из интервала
	for _, event := range manager.events[userId] {
		if (event.Date.After(since) || event.Date.Equal(since)) && event.Date.Before(to) {
			result = append(result, event)
		}
	}
	// Сортировка событий по их ID
	slices.SortFunc(result, func(a, b models.Event) int {
		return a.ID - b.ID
	})
	return result, nil
}

func NewEventManager() *EventManager {
	return &EventManager{
		events:      make(map[int]map[int]models.Event),
		nextEventID: 1,
	}
}
