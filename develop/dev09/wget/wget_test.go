package wget

import (
	"errors"
	"net/url"
	"os"
	"testing"
)

func TestWGET(t *testing.T) {
	testCases := []struct {
		name                string
		url                 string
		recursive           bool
		recursionDepth      int
		otherHosts          bool
		expectedCount       int
		expectedError       error
		expectedDirectories []string
	}{
		{
			name:                "Recursive google.com download",
			url:                 "https://www.google.com/",
			recursive:           true,
			recursionDepth:      3,
			expectedCount:       58,
			expectedError:       nil,
			expectedDirectories: []string{"./www.google.com"},
		}, {
			name:                "example.com one page",
			url:                 "https://www.example.com/",
			expectedCount:       1,
			expectedError:       nil,
			expectedDirectories: []string{"./www.example.com"},
		}, {
			name:                "Recursive example.com with other hosts",
			url:                 "https://www.example.com/",
			recursive:           true,
			recursionDepth:      1,
			otherHosts:          true,
			expectedCount:       2,
			expectedError:       nil,
			expectedDirectories: []string{"./www.example.com", "./www.iana.org"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			URL, err := url.Parse(FormatURL(testCase.url, ""))
			if err != nil {
				t.Errorf("can't run test %s: %v", testCase.name, err)
				return
			}

			options := NewOptions(URL, testCase.recursive, testCase.recursionDepth, testCase.otherHosts)
			got, err := WGET(options)

			if err != testCase.expectedError {
				t.Errorf("error: got %v, want %v", err, testCase.expectedError)
			}

			if got != testCase.expectedCount {
				t.Errorf("result: got %d, want %d", got, testCase.expectedCount)
			}

			for _, directory := range testCase.expectedDirectories {
				os.RemoveAll(directory)
			}
		})
	}
}

func TestOptionsFromArgs(t *testing.T) {
	type ExpectedOptions struct {
		address        string
		nilAddr        bool
		recursive      bool
		recursionDepth int
		otherHosts     bool
	}
	testCases := []struct {
		name            string
		arguments       []string
		expectedError   error
		expectedOptions ExpectedOptions
	}{
		{
			name:          "https://www.google.com/, Recursive, Depth 4, OtherHosts",
			arguments:     []string{"-r", "-d", "4", "-h", "https://www.google.com/"},
			expectedError: nil,
			expectedOptions: ExpectedOptions{
				address:        "https://www.google.com/",
				recursive:      true,
				recursionDepth: 4,
				otherHosts:     true,
			},
		}, {
			name:          "https://www.example.com/",
			arguments:     []string{"https://www.example.com/"},
			expectedError: nil,
			expectedOptions: ExpectedOptions{
				address: "https://www.example.com/",
			},
		}, {
			name:          "https://www.google.com/, Recursive, Depth 2",
			arguments:     []string{"-r", "-d", "2", "https://www.google.com/"},
			expectedError: nil,
			expectedOptions: ExpectedOptions{
				address:        "https://www.google.com/",
				recursive:      true,
				recursionDepth: 2,
			},
		}, {
			name:          "Unused flag",
			arguments:     []string{"-q", "https://www.google.com/"},
			expectedError: errors.New("flag provided but not defined: -q"),
			expectedOptions: ExpectedOptions{
				address:        "",
				nilAddr:        true,
				recursionDepth: 0,
			},
		}, {
			name:          "Not enough arguments",
			arguments:     []string{"-r", "-d", "2"},
			expectedError: ErrNotEnoughArguments,
			expectedOptions: ExpectedOptions{
				address:        "",
				nilAddr:        true,
				recursionDepth: 0,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			URL, err := url.Parse(FormatURL(testCase.expectedOptions.address, ""))
			if !testCase.expectedOptions.nilAddr && err != nil {
				t.Errorf("can't run test %s: %v", testCase.name, err)
				return
			}

			expectedOptions := NewOptions(URL, testCase.expectedOptions.recursive, testCase.expectedOptions.recursionDepth, testCase.expectedOptions.otherHosts)
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

			if testCase.expectedOptions.nilAddr {
				if got.BaseURL != nil {
					t.Errorf("result baseURL: got %v, want %v", got.BaseURL, nil)
				}
			} else {
				if got.BaseURL == nil {
					t.Errorf("result baseURL: got %v, want %v", got.BaseURL, URL)
				} else if *got.BaseURL != *URL {
					t.Errorf("result baseURL: got %v, want %v", *got.BaseURL, *URL)
				}
			}

			if got.Recursive != expectedOptions.Recursive {
				t.Errorf("result recursive: got %v, want %v", got.Recursive, expectedOptions.Recursive)
			}

			if got.RecursionDepth != expectedOptions.RecursionDepth {
				t.Errorf("result recursionDepth: got %v, want %v", got.RecursionDepth, expectedOptions.RecursionDepth)
			}

			if got.OtherHosts != expectedOptions.OtherHosts {
				t.Errorf("result otherHosts: got %v, want %v", got.OtherHosts, expectedOptions.OtherHosts)
			}
		})
	}
}
