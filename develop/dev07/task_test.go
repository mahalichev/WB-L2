package main

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	testCases := []struct {
		name             string
		signalDurations  []time.Duration
		expectedDuration time.Duration
	}{
		{
			name: "Default or signal",
			signalDurations: []time.Duration{
				2 * time.Hour,
				5 * time.Minute,
				1 * time.Second,
				1 * time.Hour,
				1 * time.Minute,
			},
			expectedDuration: 1 * time.Second,
		}, {
			name:             "No signal",
			signalDurations:  []time.Duration{},
			expectedDuration: 0,
		}, {
			name:             "One signal",
			signalDurations:  []time.Duration{1 * time.Second},
			expectedDuration: 1 * time.Second,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			now := time.Now()
			signals := make([]<-chan interface{}, len(testCase.signalDurations))
			for i, duration := range testCase.signalDurations {
				signals[i] = newSignal(duration)
			}
			<-or(signals...)
			got := time.Since(now)
			if got.Milliseconds() > testCase.expectedDuration.Milliseconds()+200 || got.Milliseconds() < testCase.expectedDuration.Milliseconds()-200 {
				t.Errorf("error: got %v, want %v", got, testCase.expectedDuration)
			}
		})
	}
}
