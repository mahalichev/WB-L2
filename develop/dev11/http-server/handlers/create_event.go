package handlers

import (
	"dev11/http-server/api"
	"dev11/http-server/models"
	"net/http"
	"strconv"
	"time"
)

// Обработчик создания события в календаре
func CreateEvent(manager *api.EventManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если http-метод не POST - ошибка
		if r.Method != http.MethodPost {
			ResponseJSON(w, ErrWrongRequestMethod.Error(), http.StatusBadRequest)
			return
		}

		// Парсинг данных
		if err := r.ParseForm(); err != nil {
			ResponseJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		user_id, err := strconv.Atoi(r.Form.Get("user_id"))
		if err != nil {
			ResponseJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		date, err := time.Parse(time.DateOnly, r.Form.Get("date"))
		if err != nil {
			ResponseJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		title := r.Form.Get("title")
		description := r.Form.Get("description")

		// Создание события в календаре
		event, err := manager.CreateEvent(user_id, date, title, description)
		if err != nil {
			if err == models.ErrNonPositiveID || err == models.ErrNonPositiveUserID || err == models.ErrEmptyTitle {
				ResponseJSON(w, err.Error(), http.StatusBadRequest)
				return
			}
			ResponseJSON(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		ResponseJSON(w, event, http.StatusCreated)
	}
}
