package handlers

import (
	"dev11/http-server/api"
	"dev11/http-server/models"
	"net/http"
	"strconv"
	"time"
)

// Обработчик обновления события в календаре
func UpdateEvent(manager *api.EventManager) http.HandlerFunc {
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

		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
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

		// Обновление события в календаре
		event, err := manager.UpdateEvent(id, user_id, date, title, description)
		if err != nil {
			if err == api.ErrEventNotFound || err == api.ErrUserNotFound || err == models.ErrNonPositiveID ||
				err == models.ErrNonPositiveUserID || err == models.ErrEmptyTitle {
				ResponseJSON(w, err.Error(), http.StatusBadRequest)
				return
			}
			ResponseJSON(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		ResponseJSON(w, event, http.StatusOK)
	}
}
