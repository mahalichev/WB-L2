package ntp

import (
	"testing"
	"time"
)

func TestGetNTPTime(t *testing.T) {
	testCases := []struct {
		name          string
		address       string
		expectedTime  time.Time
		expectedError string
	}{
		{
			name:          "Default NTP test",
			address:       DefaultNTPAddress,
			expectedTime:  time.Now(),
			expectedError: "",
		}, {
			name:          "Incorrect address",
			address:       "incorrect",
			expectedTime:  time.Time{},
			expectedError: "lookup incorrect: no such host",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := GetNTPTime(testCase.address)
			if (err == nil && testCase.expectedError != "") || (err != nil && err.Error() != testCase.expectedError) {
				t.Errorf("error: got %v, want %v", err, testCase.expectedError)
			}
			difference := testCase.expectedTime.Sub(got)
			if difference.Seconds() > 10 || difference.Seconds() < -10 {
				t.Errorf("result: got %s, want %s", got, testCase.expectedTime)
			}
		})
	}
}
