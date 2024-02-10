package sort

import (
	"bytes"
	"strings"
	"testing"
)

func TestSort(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		expectedOutput string
		expectedError  error
		options        Options
	}{
		{
			name:           "Default sort",
			inputText:      "test 3 test\ntest 1\ntest\ntest 2 2\n",
			expectedOutput: "test\ntest 1\ntest 2 2\ntest 3 test\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, false, false, false, false, false, false),
		}, {
			name:           "Column sort",
			inputText:      "test 3\ntest\ntest1 10 1\ntest test test\n",
			expectedOutput: "test\ntest1 10 1\ntest 3\ntest test test\n",
			expectedError:  nil,
			options:        NewOptions("", 1, false, false, false, false, false, false, false),
		}, {
			name:           "Numeric sort",
			inputText:      "2\n3\n4\n5\n6\n07\n1\n0\n",
			expectedOutput: "0\n1\n2\n3\n4\n5\n6\n07\n",
			expectedError:  nil,
			options:        NewOptions("", 0, true, false, false, false, false, false, false),
		}, {
			name:           "Month sort",
			inputText:      "February\nJan\nJanuary\nFeb\nMay\nDecember\n",
			expectedOutput: "Jan\nJanuary\nFeb\nFebruary\nMay\nDecember\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, true, false, false, false, false, false),
		}, {
			name:           "Numeric suffixes sort",
			inputText:      "3a\n1b\n3a\n12\n5f\n",
			expectedOutput: "1b\n3a\n3a\n5f\n12\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, false, true, false, false, false, false),
		}, {
			name:           "Default sort unique",
			inputText:      "test test\ntest1 test\ntest1 test1\ntest1 test\n",
			expectedOutput: "test test\ntest1 test\ntest1 test1\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, false, false, false, true, false, false),
		}, {
			name:           "Column sort unique",
			inputText:      "test test\ntest1 test\ntest1 test1\ntest1 test\n",
			expectedOutput: "test test\ntest1 test1\n",
			expectedError:  nil,
			options:        NewOptions("", 1, false, false, false, false, true, false, false),
		}, {
			name:           "Numeric sort unique",
			inputText:      "07\n1\n0001\n7\n10\n",
			expectedOutput: "1\n07\n10\n",
			expectedError:  nil,
			options:        NewOptions("", 0, true, false, false, false, true, false, false),
		}, {
			name:           "Month sort unique",
			inputText:      "May\nFeb\nJanuary\nFebruary\nJan\n",
			expectedOutput: "January\nFeb\nMay\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, true, false, false, true, false, false),
		}, {
			name:           "Numeric suffixes sort unique",
			inputText:      "3a\n2e\n3a\n4a\n2b\n1f\n",
			expectedOutput: "1f\n2b\n2e\n3a\n4a\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, false, true, false, true, false, false),
		}, {
			name:           "Check (sorted)",
			inputText:      "1f\n2b\n2e\n3a\n4a\n",
			expectedOutput: "",
			expectedError:  nil,
			options:        NewOptions("", 0, false, false, true, false, true, false, true),
		}, {
			name:           "Check (unsorted)",
			inputText:      "1f\n2b\n2e\n4a\n3a\n",
			expectedOutput: "not sorted\n",
			expectedError:  nil,
			options:        NewOptions("", 0, false, false, true, false, true, false, true),
		}, {
			name:           "Reversed column sort",
			inputText:      "test 3\ntest\ntest1 10 1\ntest test test\n",
			expectedOutput: "test test test\ntest 3\ntest1 10 1\ntest\n",
			expectedError:  nil,
			options:        NewOptions("", 1, false, false, false, true, false, false, false),
		}, {
			name:           "Column sort (ignore trailing blanks)",
			inputText:      "test4 1     \ntest2 1    \ntest5 1        \n",
			expectedOutput: "test4 1     \ntest2 1    \ntest5 1        \n",
			expectedError:  nil,
			options:        NewOptions("", 1, false, false, false, false, false, true, false),
		}, {
			name:           "Column sort (trailing blanks)",
			inputText:      "test4 1     \ntest2 1    \ntest5 1        \n",
			expectedOutput: "test2 1    \ntest4 1     \ntest5 1        \n",
			expectedError:  nil,
			options:        NewOptions("", 1, false, false, false, false, false, false, false),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.inputText)
			var buffer bytes.Buffer
			err := Sort(reader, &buffer, testCase.options)
			if testCase.expectedError != nil {
				if err == nil || err.Error() != testCase.expectedError.Error() {
					t.Errorf("error: got %v, want %v", err, testCase.expectedError)
				}
			} else {
				if err != testCase.expectedError {
					t.Errorf("error: got %v, want %v", err, testCase.expectedError)
				}
			}
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
			name:            "Only filepath",
			arguments:       []string{"./filepath.txt"},
			expectedError:   nil,
			expectedOptions: NewOptions("./filepath.txt", 0, false, false, false, false, false, false, false),
		}, {
			name:            "No arguments",
			arguments:       []string{},
			expectedError:   ErrNotEnoughArguments,
			expectedOptions: NewOptions("", 0, false, false, false, false, false, false, false),
		}, {
			name:            "Custom arguments",
			arguments:       []string{"-k", "2", "-M", "-u", "-b", "-c", "./filepath.txt"},
			expectedError:   nil,
			expectedOptions: NewOptions("./filepath.txt", 1, false, true, false, false, true, true, true),
		}, {
			name:            "Non positive column",
			arguments:       []string{"-k", "-1", "./filepath.txt"},
			expectedError:   ErrNonPositiveColumn,
			expectedOptions: NewOptions("", 0, false, false, false, false, false, false, false),
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
