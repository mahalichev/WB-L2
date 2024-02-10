package handlers

import (
	"dev11/http-server/api"
	"net/http"
	"strconv"
	"time"
)

// Обработчик получения событий пользователя из календаря
func GetEvents(manager *api.EventManager, mode string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Если http-метод не GET - ошибка
		if r.Method != http.MethodGet {
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

		since, err := time.Parse(time.DateOnly, r.Form.Get("since_date"))
		if err != nil {
			ResponseJSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		to := since
		switch mode {
		case "day":
			to = since.AddDate(0, 0, 1)
		case "week":
			to = since.AddDate(0, 0, 7)
		case "month":
			to = since.AddDate(0, 1, 0)
		}

		// Получение событий пользователя в календаре в интервале since <= date < to
		events, err := manager.GetEvents(user_id, since, to)
		if err != nil {
			if err == api.ErrEventNotFound || err == api.ErrUserNotFound {
				ResponseJSON(w, err.Error(), http.StatusBadRequest)
				return
			}
			ResponseJSON(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		ResponseJSON(w, events, http.StatusOK)
	}
}
