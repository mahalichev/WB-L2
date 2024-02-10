package handlers

import (
	"dev11/http-server/api"
	"net/http"
	"strconv"
)

// Обработчик удаления события из календаря
func DeleteEvent(manager *api.EventManager) http.HandlerFunc {
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

		// Удаление события из календаря
		err = manager.DeleteEvent(id, user_id)
		if err != nil {
			if err == api.ErrEventNotFound || err == api.ErrUserNotFound {
				ResponseJSON(w, err.Error(), http.StatusBadRequest)
				return
			}
			ResponseJSON(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		ResponseJSON(w, "ok", http.StatusOK)
	}
}
