package grep

import (
	"bytes"
	"strings"
	"testing"
)

func TestGREP(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		expectedOutput string
		options        Options
		expectedError  error
	}{
		{
			name:           "Without flags",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Grep\nGrep again\n",
			options:        NewOptions([]string{}, "Grep", 0, 0, 0, false, false, false),
			expectedError:  nil,
		}, {
			name:           "After test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Grep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\n",
			options:        NewOptions([]string{}, "Grep", 2, 0, 0, false, false, false),
			expectedError:  nil,
		}, {
			name:           "Before test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Test\nGrep\nHELLO again\nGrep again\n",
			options:        NewOptions([]string{}, "Grep", 0, 1, 0, false, false, false),
			expectedError:  nil,
		}, {
			name:           "Context test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\n",
			options:        NewOptions([]string{}, "Grep", 0, 0, 2, false, false, false),
			expectedError:  nil,
		}, {
			name:           "Context test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "2\n",
			options:        NewOptions([]string{}, "Grep", 0, 0, 0, true, false, false),
			expectedError:  nil,
		}, {
			name:           "Invert test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Hello World\nTest\nHELLO again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			options:        NewOptions([]string{}, "Grep", 0, 0, 0, false, true, false),
			expectedError:  nil,
		}, {
			name:           "Line num test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "3:Grep\n4-HELLO again\n5:Grep again\n6-Another line\n",
			options:        NewOptions([]string{}, "Grep", 1, 0, 0, false, false, true),
			expectedError:  nil,
		}, {
			name:           "Ignore case test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Grep\nGrep again\nGREP UPPER\n",
			options:        NewOptions([]string{}, "(?i)Grep", 0, 0, 0, false, false, false),
			expectedError:  nil,
		}, {
			name:           "Fixed test",
			inputText:      "Hello World\nTest\nGrep.\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Grep.\n",
			options:        NewOptions([]string{}, "Grep\\.", 0, 0, 0, false, false, false),
			expectedError:  nil,
		}, {
			name:           "Pattern test",
			inputText:      "Hello World\nTest\nGrep\nHELLO again\nGrep again\nAnother line\nGREP UPPER\nSample text here\nhello one more time\nEnd of file\n",
			expectedOutput: "Grep again\n",
			options:        NewOptions([]string{}, "Grep.", 0, 0, 0, false, false, false),
			expectedError:  nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.inputText)
			var buffer bytes.Buffer
			err := GREP(reader, &buffer, testCase.options)
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
				t.Errorf("result: got %s, want %s", got, testCase.expectedOutput)
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
			name:            "All flags",
			arguments:       []string{"-A", "3", "-B", "1", "-C", "2", "-c", "-i", "-v", "-F", "-n", "test.", "test.txt"},
			expectedError:   nil,
			expectedOptions: NewOptions([]string{"test.txt"}, "(?i)test\\.", 3, 1, 2, true, true, true),
		}, {
			name:            "No arguments",
			arguments:       []string{},
			expectedError:   ErrNotEnoughArguments,
			expectedOptions: NewOptions([]string{}, "", 0, 0, 0, false, false, false),
		}, {
			name:            "Custom flags",
			arguments:       []string{"-C", "2", "-i", "-v", "-c", "test."},
			expectedError:   nil,
			expectedOptions: NewOptions([]string{}, "(?i)test.", 0, 0, 2, true, true, false),
		}, {
			name:            "No pattern",
			arguments:       []string{"-C", "2", "-i", "-v", "-c"},
			expectedError:   ErrNotEnoughArguments,
			expectedOptions: NewOptions([]string{}, "", 0, 0, 0, false, false, false),
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
			if len(got.Filepaths) != len(testCase.expectedOptions.Filepaths) {
				t.Errorf("result filepaths: got %v, want %v", got, testCase.expectedOptions.Filepaths)
			} else {
				for i, file := range testCase.expectedOptions.Filepaths {
					if got.Filepaths[i] != file {
						t.Errorf("result filepaths: got %v, want %v", got.Filepaths[i], file)
						break
					}
				}
			}
			if got.After != testCase.expectedOptions.After {
				t.Errorf("result after: got %v, want %v", got.After, testCase.expectedOptions.After)
			}
			if got.Before != testCase.expectedOptions.Before {
				t.Errorf("result before: got %v, want %v", got.Before, testCase.expectedOptions.Before)
			}
			if got.Context != testCase.expectedOptions.Context {
				t.Errorf("result context: got %v, want %v", got.Context, testCase.expectedOptions.Context)
			}
			if got.Count != testCase.expectedOptions.Count {
				t.Errorf("result count: got %v, want %v", got.Count, testCase.expectedOptions.Count)
			}
			if got.Invert != testCase.expectedOptions.Invert {
				t.Errorf("result invert: got %v, want %v", got.Invert, testCase.expectedOptions.Invert)
			}
			if got.LineNum != testCase.expectedOptions.LineNum {
				t.Errorf("result line num: got %v, want %v", got.LineNum, testCase.expectedOptions.LineNum)
			}
			if got.Pattern != testCase.expectedOptions.Pattern {
				t.Errorf("result pattern: got %v, want %v", got.Pattern, testCase.expectedOptions.Pattern)
			}
		})
	}
}
