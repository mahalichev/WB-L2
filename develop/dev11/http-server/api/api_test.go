package api

import (
	"dev11/http-server/models"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	testCases := []struct {
		name          string
		initState     map[int]map[int]models.Event
		initNextID    int
		userId        int
		date          time.Time
		title         string
		description   string
		expectedState map[int]map[int]models.Event
		expectedError error
	}{
		{
			name:        "Default create",
			initState:   map[int]map[int]models.Event{},
			initNextID:  1,
			userId:      1,
			date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			title:       "TestTitle",
			description: "TestDescription",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: nil,
		}, {
			name: "Another user",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  2,
			userId:      2,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitle2",
			description: "TestDescription2",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
				2: {
					2: {
						Date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      2,
						Title:       "TestTitle2",
						Description: "TestDescription2",
					},
				},
			},
			expectedError: nil,
		}, {
			name:          "Non positive user_id",
			initState:     map[int]map[int]models.Event{},
			initNextID:    1,
			userId:        0,
			date:          time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			title:         "TestTitle",
			description:   "TestDescription",
			expectedState: map[int]map[int]models.Event{},
			expectedError: models.ErrNonPositiveUserID,
		}, {
			name:          "Empty title",
			initState:     map[int]map[int]models.Event{},
			initNextID:    1,
			userId:        1,
			date:          time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			title:         "",
			description:   "TestDescription",
			expectedState: map[int]map[int]models.Event{},
			expectedError: models.ErrEmptyTitle,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			manager := NewEventManager()
			manager.events = testCase.initState
			manager.nextEventID = testCase.initNextID
			_, err := manager.CreateEvent(testCase.userId, testCase.date, testCase.title, testCase.description)
			if err != testCase.expectedError {
				t.Errorf("error: got %v, expected %v", err, testCase.expectedError)
			}
			if len(testCase.expectedState) != len(manager.events) {
				t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
				return
			}
			for userId, userEvents := range testCase.expectedState {
				if _, ok := manager.events[userId]; !ok {
					t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
					return
				}
				if len(userEvents) != len(manager.events[userId]) {
					t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
					return
				}
				for eventId, expectedEvent := range userEvents {
					if event, ok := manager.events[userId][eventId]; !ok {
						t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
						return
					} else if event != expectedEvent {
						t.Errorf("result: got %v, expected %v", event, expectedEvent)
						return
					}
				}
			}
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	testCases := []struct {
		name          string
		initState     map[int]map[int]models.Event
		initNextID    int
		eventId       int
		userId        int
		expectedState map[int]map[int]models.Event
		expectedError error
	}{
		{
			name: "Default delete",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:    2,
			eventId:       1,
			userId:        1,
			expectedState: map[int]map[int]models.Event{1: {}},
			expectedError: nil,
		}, {
			name: "Delete from user",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle1",
						Description: "TestDescription1",
					},
				},
				2: {
					2: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      2,
						Title:       "TestTitle2",
						Description: "TestDescription2",
					},
				},
			},
			initNextID: 3,
			eventId:    2,
			userId:     2,
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle1",
						Description: "TestDescription1",
					},
				},
				2: {},
			},
			expectedError: nil,
		}, {
			name: "User not found",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID: 2,
			eventId:    2,
			userId:     2,
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: ErrUserNotFound,
		}, {
			name: "Event not found",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID: 2,
			eventId:    2,
			userId:     1,
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: ErrEventNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			manager := NewEventManager()
			manager.events = testCase.initState
			manager.nextEventID = testCase.initNextID
			err := manager.DeleteEvent(testCase.eventId, testCase.userId)
			if err != testCase.expectedError {
				t.Errorf("error: got %v, expected %v", err, testCase.expectedError)
			}
			if len(testCase.expectedState) != len(manager.events) {
				t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
				return
			}
			for userId, userEvents := range testCase.expectedState {
				if _, ok := manager.events[userId]; !ok {
					t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
					return
				}
				if len(userEvents) != len(manager.events[userId]) {
					t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
					return
				}
				for eventId, expectedEvent := range userEvents {
					if event, ok := manager.events[userId][eventId]; !ok {
						t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
						return
					} else if event != expectedEvent {
						t.Errorf("result: got %v, expected %v", event, expectedEvent)
						return
					}
				}
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	testCases := []struct {
		name          string
		initState     map[int]map[int]models.Event
		initNextID    int
		eventId       int
		userId        int
		date          time.Time
		title         string
		description   string
		expectedState map[int]map[int]models.Event
		expectedError error
	}{
		{
			name: "Default update",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     1,
			userId:      1,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitleUpdated",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitleUpdated",
						Description: "TestDescriptionUpdated",
					},
				},
			},
			expectedError: nil,
		}, {
			name: "Another user",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
				2: {
					2: {
						Date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      2,
						Title:       "TestTitle2",
						Description: "TestDescription2",
					},
				},
			},
			initNextID:  2,
			eventId:     2,
			userId:      2,
			date:        time.Date(2006, 10, 12, 0, 0, 0, 0, time.UTC),
			title:       "TestTitle3",
			description: "TestDescription3",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
				2: {
					2: {
						Date:        time.Date(2006, 10, 12, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      2,
						Title:       "TestTitle3",
						Description: "TestDescription3",
					},
				},
			},
			expectedError: nil,
		}, {
			name: "Non positive user_id",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     1,
			userId:      0,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitleUpdated",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: models.ErrNonPositiveUserID,
		}, {
			name: "Non positive event_id",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     0,
			userId:      1,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitleUpdated",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: models.ErrNonPositiveID,
		}, {
			name: "Empty title",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     1,
			userId:      1,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: models.ErrEmptyTitle,
		}, {
			name: "Non positive event_id",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     0,
			userId:      1,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitleUpdated",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: models.ErrNonPositiveID,
		}, {
			name: "User not found",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     1,
			userId:      2,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitleUpdated",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: ErrUserNotFound,
		}, {
			name: "User not found",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			initNextID:  1,
			eventId:     2,
			userId:      1,
			date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			title:       "TestTitleUpdated",
			description: "TestDescriptionUpdated",
			expectedState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle",
						Description: "TestDescription",
					},
				},
			},
			expectedError: ErrEventNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			manager := NewEventManager()
			manager.events = testCase.initState
			manager.nextEventID = testCase.initNextID
			_, err := manager.UpdateEvent(testCase.eventId, testCase.userId, testCase.date, testCase.title, testCase.description)
			if err != testCase.expectedError {
				t.Errorf("error: got %v, expected %v", err, testCase.expectedError)
			}
			if len(testCase.expectedState) != len(manager.events) {
				t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
				return
			}
			for userId, userEvents := range testCase.expectedState {
				if _, ok := manager.events[userId]; !ok {
					t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
					return
				}
				if len(userEvents) != len(manager.events[userId]) {
					t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
					return
				}
				for eventId, expectedEvent := range userEvents {
					if event, ok := manager.events[userId][eventId]; !ok {
						t.Errorf("result: got %v, expected %v", manager.events, testCase.expectedState)
						return
					} else if event != expectedEvent {
						t.Errorf("result: got %v, expected %v", event, expectedEvent)
						return
					}
				}
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	testCases := []struct {
		name           string
		initState      map[int]map[int]models.Event
		initNextID     int
		userId         int
		since_date     time.Time
		to             time.Time
		expectedResult []models.Event
		expectedError  error
	}{
		{
			name: "Default get (day)",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle1",
						Description: "TestDescription1",
					},
					2: {
						Date:        time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      1,
						Title:       "TestTitle2",
						Description: "TestDescription2",
					},
				},
				2: {
					3: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          3,
						UserID:      2,
						Title:       "TestTitle3",
						Description: "TestDescription3",
					},
				},
			},
			initNextID:     4,
			userId:         1,
			since_date:     time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			to:             time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			expectedResult: []models.Event{models.NewEvent(1, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1")},
			expectedError:  nil,
		}, {
			name: "Default get (week)",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle1",
						Description: "TestDescription1",
					},
					2: {
						Date:        time.Date(2006, 10, 16, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      1,
						Title:       "TestTitle2",
						Description: "TestDescription2",
					},
					3: {
						Date:        time.Date(2006, 10, 17, 0, 0, 0, 0, time.UTC),
						ID:          3,
						UserID:      1,
						Title:       "TestTitle3",
						Description: "TestDescription3",
					},
				},
				2: {
					4: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          4,
						UserID:      2,
						Title:       "TestTitle3",
						Description: "TestDescription3",
					},
				},
			},
			initNextID: 5,
			userId:     1,
			since_date: time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			to:         time.Date(2006, 10, 17, 0, 0, 0, 0, time.UTC),
			expectedResult: []models.Event{
				models.NewEvent(1, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1"),
				models.NewEvent(2, 1, time.Date(2006, 10, 16, 0, 0, 0, 0, time.UTC), "TestTitle2", "TestDescription2"),
			},
			expectedError: nil,
		}, {
			name: "Default get (month)",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle1",
						Description: "TestDescription1",
					},
					2: {
						Date:        time.Date(2006, 11, 9, 0, 0, 0, 0, time.UTC),
						ID:          2,
						UserID:      1,
						Title:       "TestTitle2",
						Description: "TestDescription2",
					},
					3: {
						Date:        time.Date(2006, 11, 10, 0, 0, 0, 0, time.UTC),
						ID:          3,
						UserID:      1,
						Title:       "TestTitle3",
						Description: "TestDescription3",
					},
				},
				2: {
					4: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          4,
						UserID:      2,
						Title:       "TestTitle3",
						Description: "TestDescription3",
					},
				},
			},
			initNextID: 5,
			userId:     1,
			since_date: time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			to:         time.Date(2006, 11, 10, 0, 0, 0, 0, time.UTC),
			expectedResult: []models.Event{
				models.NewEvent(1, 1, time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC), "TestTitle1", "TestDescription1"),
				models.NewEvent(2, 1, time.Date(2006, 11, 9, 0, 0, 0, 0, time.UTC), "TestTitle2", "TestDescription2"),
			},
			expectedError: nil,
		}, {
			name: "User not found",
			initState: map[int]map[int]models.Event{
				1: {
					1: {
						Date:        time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
						ID:          1,
						UserID:      1,
						Title:       "TestTitle1",
						Description: "TestDescription1",
					},
				},
			},
			initNextID:     2,
			userId:         2,
			since_date:     time.Date(2006, 10, 10, 0, 0, 0, 0, time.UTC),
			to:             time.Date(2006, 10, 11, 0, 0, 0, 0, time.UTC),
			expectedResult: []models.Event{},
			expectedError:  ErrUserNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			manager := NewEventManager()
			manager.events = testCase.initState
			manager.nextEventID = testCase.initNextID
			got, err := manager.GetEvents(testCase.userId, testCase.since_date, testCase.to)
			if err != testCase.expectedError {
				t.Errorf("error: got %v, expected %v", err, testCase.expectedError)
			}
			if len(testCase.expectedResult) != len(got) {
				t.Errorf("result: got %v, expected %v", got, testCase.expectedResult)
				return
			}
			for id, event := range testCase.expectedResult {
				if got[id] != event {
					t.Errorf("result: got %v, expected %v", got[id], event)
					return
				}
			}
		})
	}
}
