package unpack

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	testCases := []struct {
		name           string
		packed         string
		expectedString string
		expectedError  error
	}{
		{
			name:           "Default packed string",
			packed:         `a4bc2d5e`,
			expectedString: `aaaabccddddde`,
			expectedError:  nil,
		}, {
			name:           "Only unique characters",
			packed:         `abcd`,
			expectedString: `abcd`,
			expectedError:  nil,
		}, {
			name:           "Only count",
			packed:         `45`,
			expectedString: ``,
			expectedError:  ErrIncorrectPackedString,
		}, {
			name:           "Empty string",
			packed:         ``,
			expectedString: ``,
			expectedError:  nil,
		}, {
			name:           "Escape digits",
			packed:         `qwe\4\5`,
			expectedString: "qwe45",
			expectedError:  nil,
		}, {
			name:           "Unpack escaping digit",
			packed:         `qwe\45`,
			expectedString: `qwe44444`,
			expectedError:  nil,
		}, {
			name:           "Unpack backslash",
			packed:         `qwe\\5`,
			expectedString: `qwe\\\\\`,
			expectedError:  nil,
		}, {
			name:           "Unpack letter",
			packed:         `\a`,
			expectedString: ``,
			expectedError:  ErrIncorrectPackedString,
		}, {
			name:           "Zero count",
			packed:         `a0`,
			expectedString: ``,
			expectedError:  ErrIncorrectPackedString,
		}, {
			name:           "One backslash",
			packed:         `\`,
			expectedString: ``,
			expectedError:  ErrIncorrectPackedString,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := UnpackString(testCase.packed)
			if err != testCase.expectedError {
				t.Errorf("error: got %v, want %v", err, testCase.expectedError)
			}
			if got != testCase.expectedString {
				t.Errorf("result: got %s, want %s", got, testCase.expectedString)
			}
		})
	}
}

func TestRepeatLetter(t *testing.T) {
	testCases := []struct {
		name     string
		letter   string
		count    int
		expected string
	}{
		{
			name:     "Backslash",
			letter:   `\`,
			count:    5,
			expected: `\\\\\`,
		}, {
			name:     "Digit",
			letter:   `3`,
			count:    10,
			expected: `3333333333`,
		}, {
			name:     "Zero count",
			letter:   "a",
			count:    0,
			expected: "a",
		}, {
			name:     "Empty letter",
			letter:   "",
			count:    1,
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := RepeatLetter(testCase.letter, testCase.count)
			if got != testCase.expected {
				t.Errorf("got %v, want %v", got, testCase.expected)
			}
		})
	}
}
