package handlers

import (
	"dev11/http-server/api"
	"dev11/http-server/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateEventHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		form           string
		expectedOutput string
		expectedStatus int
	}{
		{
			name:           "Default request",
			method:         "POST",
			form:           "user_id=1&date=2006-10-10&title=TestTitle&description=TestDescription",
			expectedOutput: "{\"result\":{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle\",\"Description\":\"TestDescription\",\"ID\":1,\"UserID\":1}}\n",
			expectedStatus: http.StatusCreated,
		}, {
			name:           "Wrong user_id",
			method:         "POST",
			form:           "user_id=Test&date=2006-10-10&title=TestTitle&description=TestDescription",
			expectedOutput: "{\"error\":\"strconv.Atoi: parsing \\\"Test\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong date format",
			method:         "POST",
			form:           "user_id=1&date=2006&title=TestTitle&description=TestDescription",
			expectedOutput: "{\"error\":\"parsing time \\\"2006\\\" as \\\"2006-01-02\\\": cannot parse \\\"\\\" as \\\"-\\\"\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "No title",
			method:         "POST",
			form:           "user_id=1&date=2006-10-10&description=TestDescription",
			expectedOutput: "{\"error\":\"title cannot be empty\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "No description",
			method:         "POST",
			form:           "user_id=1&date=2006-10-10&title=TestTitle",
			expectedOutput: "{\"result\":{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle\",\"Description\":\"\",\"ID\":1,\"UserID\":1}}\n",
			expectedStatus: http.StatusCreated,
		}, {
			name:           "Wrong method",
			method:         "GET",
			form:           "user_id=1&date=2006-10-10&title=TestTitle&description=TestDescription",
			expectedOutput: "{\"error\":\"wrong request method\"}\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(testCase.method, "localhost:3000/create_event", strings.NewReader(testCase.form))
			if err != nil {
				t.Error(err)
				return
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateEvent(api.NewEventManager()))
			handler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != testCase.expectedStatus {
				t.Errorf("status: got %v want %v", responseRecorder.Code, testCase.expectedStatus)
			}
			if responseRecorder.Body.String() != testCase.expectedOutput {
				t.Errorf("result: got %v want %v", responseRecorder.Body.String(), testCase.expectedOutput)
			}
		})
	}
}

func TestUpdateEventHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		form           string
		managerEntries []models.Event
		expectedOutput string
		expectedStatus int
	}{
		{
			name:           "Default request",
			method:         "POST",
			form:           "id=1&user_id=1&date=2006-10-11&title=NewTitle&description=NewDescription",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"result\":{\"Date\":\"2006-10-11T00:00:00Z\",\"Title\":\"NewTitle\",\"Description\":\"NewDescription\",\"ID\":1,\"UserID\":1}}\n",
			expectedStatus: http.StatusOK,
		}, {
			name:           "Wrong id",
			method:         "POST",
			form:           "id=Test&user_id=1&date=2006-10-11&title=NewTitle&description=NewDescription",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"strconv.Atoi: parsing \\\"Test\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong user_id",
			method:         "POST",
			form:           "id=1&user_id=Test&date=2006-10-11&title=NewTitle&description=NewDescription",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"strconv.Atoi: parsing \\\"Test\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong date format",
			method:         "POST",
			form:           "id=1&user_id=1&date=2006&title=NewTitle&description=NewDescription",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"parsing time \\\"2006\\\" as \\\"2006-01-02\\\": cannot parse \\\"\\\" as \\\"-\\\"\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "No title",
			method:         "POST",
			form:           "id=1&user_id=1&date=2006-10-10&description=TestDescription",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"title cannot be empty\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "No description",
			method:         "POST",
			form:           "id=1&user_id=1&date=2006-10-10&title=TestTitle",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"result\":{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle\",\"Description\":\"\",\"ID\":1,\"UserID\":1}}\n",
			expectedStatus: http.StatusOK,
		}, {
			name:           "Event not found",
			method:         "POST",
			form:           "id=2&user_id=1&date=2006-10-10&title=TestTitle",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"event with given ID not found\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "User not found",
			method:         "POST",
			form:           "id=1&user_id=2&date=2006-10-10&title=TestTitle",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"user with given ID not found\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong method",
			method:         "GET",
			form:           "id=1&user_id=1&date=2006-10-11&title=NewTitle&description=NewDescription",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"wrong request method\"}\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(testCase.method, "localhost:3000/update_event", strings.NewReader(testCase.form))
			if err != nil {
				t.Error(err)
				return
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			responseRecorder := httptest.NewRecorder()
			manager := api.NewEventManager()
			for _, event := range testCase.managerEntries {
				manager.CreateEvent(event.UserID, event.Date, event.Title, event.Description)
			}
			handler := http.HandlerFunc(UpdateEvent(manager))
			handler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != testCase.expectedStatus {
				t.Errorf("status: got %v want %v", responseRecorder.Code, testCase.expectedStatus)
			}
			if responseRecorder.Body.String() != testCase.expectedOutput {
				t.Errorf("result: got %v want %v", responseRecorder.Body.String(), testCase.expectedOutput)
			}
		})
	}
}

func TestDeleteEventHandler(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		form           string
		managerEntries []models.Event
		expectedOutput string
		expectedStatus int
	}{
		{
			name:           "Default request",
			method:         "POST",
			form:           "id=1&user_id=1",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"result\":\"ok\"}\n",
			expectedStatus: http.StatusOK,
		}, {
			name:           "Wrong id",
			method:         "POST",
			form:           "id=Test&user_id=1",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"strconv.Atoi: parsing \\\"Test\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong user_id",
			method:         "POST",
			form:           "id=1&user_id=Test",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"strconv.Atoi: parsing \\\"Test\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Event not found",
			method:         "POST",
			form:           "id=2&user_id=1",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"event with given ID not found\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "User not found",
			method:         "POST",
			form:           "id=1&user_id=2",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"user with given ID not found\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong method",
			method:         "GET",
			form:           "id=1&user_id=1",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"wrong request method\"}\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(testCase.method, "localhost:3000/delete_event", strings.NewReader(testCase.form))
			if err != nil {
				t.Error(err)
				return
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			responseRecorder := httptest.NewRecorder()
			manager := api.NewEventManager()
			for _, event := range testCase.managerEntries {
				manager.CreateEvent(event.UserID, event.Date, event.Title, event.Description)
			}
			handler := http.HandlerFunc(DeleteEvent(manager))
			handler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != testCase.expectedStatus {
				t.Errorf("status: got %v want %v", responseRecorder.Code, testCase.expectedStatus)
			}
			if responseRecorder.Body.String() != testCase.expectedOutput {
				t.Errorf("result: got %v want %v", responseRecorder.Body.String(), testCase.expectedOutput)
			}
		})
	}
}

func TestGetEventsHandler(t *testing.T) {
	testCases := []struct {
		name           string
		url            string
		method         string
		mode           string
		managerEntries []models.Event
		expectedOutput string
		expectedStatus int
	}{
		{
			name:   "One day request",
			url:    "localhost:3000/events_for_day?user_id=1&since_date=2006-10-10",
			method: "GET",
			mode:   "day",
			managerEntries: []models.Event{
				models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1"),
				models.NewEvent(0, 1, time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC), "TestTitle2", "TestDescription2"),
				models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle3", "TestDescription3"),
				models.NewEvent(0, 2, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription"),
			},
			expectedOutput: "{\"result\":[{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle1\",\"Description\":\"TestDescription1\",\"ID\":1,\"UserID\":1},{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle3\",\"Description\":\"TestDescription3\",\"ID\":3,\"UserID\":1}]}\n",
			expectedStatus: http.StatusOK,
		}, {
			name:           "Wrong date format",
			url:            "localhost:3000/events_for_day?user_id=1&since_date=2006",
			method:         "GET",
			mode:           "day",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"parsing time \\\"2006\\\" as \\\"2006-01-02\\\": cannot parse \\\"\\\" as \\\"-\\\"\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "Wrong user_id",
			url:            "localhost:3000/events_for_day?user_id=Test&since_date=2006-10-10",
			method:         "GET",
			mode:           "day",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"strconv.Atoi: parsing \\\"Test\\\": invalid syntax\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:           "User not found",
			url:            "localhost:3000/events_for_day?user_id=2&since_date=2006-10-10",
			method:         "GET",
			mode:           "day",
			managerEntries: []models.Event{models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription")},
			expectedOutput: "{\"error\":\"user with given ID not found\"}\n",
			expectedStatus: http.StatusBadRequest,
		}, {
			name:   "One week request",
			url:    "localhost:3000/events_for_week?user_id=1&since_date=2006-10-10",
			method: "GET",
			mode:   "week",
			managerEntries: []models.Event{
				models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1"),
				models.NewEvent(0, 1, time.Date(2006, 10, 17, 0, 0, 0, 0, time.UTC), "TestTitle2", "TestDescription2"),
				models.NewEvent(0, 1, time.Date(2006, 10, 16, 0, 0, 0, 0, time.UTC), "TestTitle3", "TestDescription3"),
				models.NewEvent(0, 2, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription"),
			},
			expectedOutput: "{\"result\":[{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle1\",\"Description\":\"TestDescription1\",\"ID\":1,\"UserID\":1},{\"Date\":\"2006-10-16T00:00:00Z\",\"Title\":\"TestTitle3\",\"Description\":\"TestDescription3\",\"ID\":3,\"UserID\":1}]}\n",
			expectedStatus: http.StatusOK,
		}, {
			name:   "One month request",
			url:    "localhost:3000/events_for_month?user_id=1&since_date=2006-10-10",
			method: "GET",
			mode:   "month",
			managerEntries: []models.Event{
				models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1"),
				models.NewEvent(0, 1, time.Date(2006, 11, 10, 0, 0, 0, 0, time.UTC), "TestTitle2", "TestDescription2"),
				models.NewEvent(0, 1, time.Date(2006, 11, 9, 0, 0, 0, 0, time.UTC), "TestTitle3", "TestDescription3"),
				models.NewEvent(0, 2, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription"),
			},
			expectedOutput: "{\"result\":[{\"Date\":\"2006-10-10T00:00:00Z\",\"Title\":\"TestTitle1\",\"Description\":\"TestDescription1\",\"ID\":1,\"UserID\":1},{\"Date\":\"2006-11-09T00:00:00Z\",\"Title\":\"TestTitle3\",\"Description\":\"TestDescription3\",\"ID\":3,\"UserID\":1}]}\n",
			expectedStatus: http.StatusOK,
		}, {
			name:   "Wrong method",
			url:    "localhost:3000/events_for_day?user_id=1&since_date=2006-10-10",
			method: "POST",
			mode:   "day",
			managerEntries: []models.Event{
				models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1"),
				models.NewEvent(0, 1, time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC), "TestTitle2", "TestDescription2"),
				models.NewEvent(0, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle3", "TestDescription3"),
				models.NewEvent(0, 2, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle", "TestDescription"),
			},
			expectedOutput: "{\"error\":\"wrong request method\"}\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(testCase.method, testCase.url, nil)
			if err != nil {
				t.Error(err)
				return
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			responseRecorder := httptest.NewRecorder()
			manager := api.NewEventManager()
			for _, event := range testCase.managerEntries {
				manager.CreateEvent(event.UserID, event.Date, event.Title, event.Description)
			}
			handler := http.HandlerFunc(GetEvents(manager, testCase.mode))
			handler.ServeHTTP(responseRecorder, request)
			if responseRecorder.Code != testCase.expectedStatus {
				t.Errorf("status: got %v want %v", responseRecorder.Code, testCase.expectedStatus)
			}
			if responseRecorder.Body.String() != testCase.expectedOutput {
				t.Errorf("result: got %v want %v", responseRecorder.Body.String(), testCase.expectedOutput)
			}
		})
	}
}
