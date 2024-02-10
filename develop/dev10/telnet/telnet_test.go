package telnet

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestTelnet(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		expectedOutput string
		options        Options
	}{
		{
			name:           "One string",
			inputText:      "test\n",
			expectedOutput: "server: test\n",
			options:        NewOptions("127.0.0.1", "3010", time.Duration(10*time.Second)),
		}, {
			name:           "Three strings",
			inputText:      "test1\ntest2\ntest3\n",
			expectedOutput: "server: test1\nserver: test2\nserver: test3\n",
			options:        NewOptions("127.0.0.1", "3010", time.Duration(10*time.Second)),
		}, {
			name:           "No input",
			inputText:      "",
			expectedOutput: "",
			options:        NewOptions("127.0.0.1", "3010", time.Duration(10*time.Second)),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.inputText)
			var buffer bytes.Buffer
			Telnet(reader, &buffer, testCase.options)
			if got := buffer.String(); got != testCase.expectedOutput {
				t.Errorf("got %s, want %s", got, testCase.expectedOutput)
			}
		})
	}
}

func TestParseArguments(t *testing.T) {
	testCases := []struct {
		name            string
		arguments       []string
		expectedError   error
		expectedOptions Options
	}{
		{
			name:            "Only host and port",
			arguments:       []string{"google.com", "3000"},
			expectedError:   nil,
			expectedOptions: NewOptions("google.com", "3000", time.Duration(10*time.Second)),
		}, {
			name:            "No arguments",
			arguments:       []string{},
			expectedError:   ErrNotEnoughArguments,
			expectedOptions: NewOptions("", "", time.Duration(0)),
		}, {
			name:            "Custom timeout 2h",
			arguments:       []string{"--timeout=2h", "google.com", "3000"},
			expectedError:   nil,
			expectedOptions: NewOptions("google.com", "3000", time.Duration(2*time.Hour)),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ParseArguments(testCase.arguments)

			if testCase.expectedError != nil {
				if err == nil || err.Error() != testCase.expectedError.Error() {
					t.Errorf("error: got %v, want %v", err, testCase.expectedError)
				}
			} else {
				if err != testCase.expectedError {
					t.Errorf("error: got %v, want %v", err, testCase.expectedError)
				}
			}

			if got != testCase.expectedOptions {
				t.Errorf("result recursive: got %v, want %v", got, testCase.expectedOptions)
			}
		})
	}
}
