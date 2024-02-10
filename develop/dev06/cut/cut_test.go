package cut

import (
	"bytes"
	"strings"
	"testing"
)

func CompareFieldsParams(t *testing.T, got, expected FieldsParams) {
	if got.FromStartTo != expected.FromStartTo {
		t.Errorf("result: FromStartTo got %v, want %v", got, expected)
	}

	if got.IsFromStart != expected.IsFromStart {
		t.Errorf("result: IsFromStart got %v, want %v", got.IsFromStart, expected.IsFromStart)
	}

	if got.FromToEnd != expected.FromToEnd {
		t.Errorf("result: FromToEnd got %v, want %v", got.FromToEnd, expected.FromToEnd)
	}

	if got.IsToEnd != expected.IsToEnd {
		t.Errorf("result: IsToEnd got %v, want %v", got.IsToEnd, expected.IsToEnd)
	}

	if len(got.Fields) != len(expected.Fields) {
		t.Errorf("result: Fields got %v, want %v", got.Fields, expected.Fields)
		return
	}
	for i := range got.Fields {
		if got.Fields[i] != expected.Fields[i] {
			t.Errorf("result: Fields got %v, want %v", got.Fields, expected.Fields)
			return
		}
	}
}

func TestGetFieldsFromString(t *testing.T) {
	testCases := []struct {
		name           string
		fieldsString   string
		expectedOutput FieldsParams
		expectedError  error
	}{
		{
			name:           "Only commas",
			fieldsString:   "1,3,6,10",
			expectedOutput: FieldsParams{Fields: []int{0, 2, 5, 9}},
			expectedError:  nil,
		}, {
			name:           "Only intervals",
			fieldsString:   "1-3,6-10",
			expectedOutput: FieldsParams{Fields: []int{0, 1, 2, 5, 6, 7, 8, 9}},
			expectedError:  nil,
		}, {
			name:           "FromStartTo",
			fieldsString:   "-3,5",
			expectedOutput: FieldsParams{Fields: []int{4}, FromStartTo: 2, IsFromStart: true},
			expectedError:  nil,
		}, {
			name:           "FromToEnd",
			fieldsString:   "1,6-",
			expectedOutput: FieldsParams{Fields: []int{0}, FromToEnd: 5, IsToEnd: true},
			expectedError:  nil,
		}, {
			name:           "Combined",
			fieldsString:   "-3,5-8,4,11,20-",
			expectedOutput: FieldsParams{Fields: []int{4, 5, 6, 7, 3, 10}, FromStartTo: 2, IsFromStart: true, FromToEnd: 19, IsToEnd: true},
			expectedError:  nil,
		}, {
			name:           "Combined crossed",
			fieldsString:   "-3,5-8,6,5,21,20-",
			expectedOutput: FieldsParams{Fields: []int{4, 5, 6, 7}, FromStartTo: 2, IsFromStart: true, FromToEnd: 19, IsToEnd: true},
			expectedError:  nil,
		}, {
			name:           "Wrong input",
			fieldsString:   "-",
			expectedOutput: FieldsParams{},
			expectedError:  ErrWrongFieldEntry,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := GetFieldsFromString(testCase.fieldsString)

			if err != testCase.expectedError {
				t.Errorf("error: got %v, want %v", err, testCase.expectedError)
			}

			CompareFieldsParams(t, got, testCase.expectedOutput)
		})
	}
}

func TestParseArguments(t *testing.T) {
	testCases := []struct {
		name           string
		arguments      []string
		expectedOutput Options
		expectedError  error
	}{
		{
			name:           "All flags",
			arguments:      []string{"-d", " ", "-s", "-f", "1,3,5"},
			expectedOutput: NewOptions(FieldsParams{Fields: []int{0, 2, 4}}, " ", true),
			expectedError:  nil,
		}, {
			name:           "Default delimiter",
			arguments:      []string{"-s", "-f", "1,3,5"},
			expectedOutput: NewOptions(FieldsParams{Fields: []int{0, 2, 4}}, "\t", true),
			expectedError:  nil,
		}, {
			name:           "With no separated strings",
			arguments:      []string{"-f", "1-5", "-d", " "},
			expectedOutput: NewOptions(FieldsParams{Fields: []int{0, 1, 2, 3, 4}}, " ", false),
			expectedError:  nil,
		}, {
			name:           "No fields flag",
			arguments:      []string{"-d", " ", "-s"},
			expectedOutput: Options{},
			expectedError:  ErrWrongFieldEntry,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ParseArguments(testCase.arguments)

			if err != testCase.expectedError {
				t.Errorf("error: got %v, want %v", err, testCase.expectedError)
			}

			if got.Delimiter != testCase.expectedOutput.Delimiter {
				t.Errorf("result: Delimiter got %v, want %v", got.Delimiter, testCase.expectedOutput.Delimiter)
			}

			if got.Separated != testCase.expectedOutput.Separated {
				t.Errorf("result: Separated got %v, want %v", got.Separated, testCase.expectedOutput.Separated)
			}

			CompareFieldsParams(t, got.FieldsParams, testCase.expectedOutput.FieldsParams)
		})
	}
}

func TestCut(t *testing.T) {
	testCases := []struct {
		name           string
		inputText      string
		expectedOutput string
		expectedError  error
		options        Options
	}{
		{
			name:           "Default cut",
			inputText:      "test11\ttest12\ttest13\ttest14\ttest15\ntest21\ttest22\ttest23\ttest24\ttest25\ntest31\ttest32\ttest33\ttest34\ttest35\ntest41 test42 test43 test44 test45\n",
			expectedOutput: "test11\ttest13\ttest15\t\ntest21\ttest23\ttest25\t\ntest31\ttest33\ttest35\t\ntest41 test42 test43 test44 test45\n",
			expectedError:  nil,
			options:        NewOptions(FieldsParams{Fields: []int{0, 2, 4}}, "\t", false),
		}, {
			name:           "Custom separator",
			inputText:      "test11 test12 test13 test14 test15\ntest21\ttest22\ttest23\ttest24\ttest25\ntest31\ttest32\ttest33\ttest34\ttest35\ntest41 test42 test43 test44 test45\n",
			expectedOutput: "test12 test14 \ntest42 test44 \n",
			expectedError:  nil,
			options:        NewOptions(FieldsParams{Fields: []int{1, 3}}, " ", true),
		}, {
			name:           "From start to",
			inputText:      "test11\ttest12\ttest13\ttest14\ntest21\ttest22\ttest23\ntest31\ttest32\ttest33\ttest34\ttest35\ntest41 test42\n",
			expectedOutput: "test11\ttest12\ttest13\t\ntest21\ttest22\ttest23\t\ntest31\ttest32\ttest33\t\ntest41 test42\n",
			expectedError:  nil,
			options:        NewOptions(FieldsParams{IsFromStart: true, FromStartTo: 2}, "\t", false),
		}, {
			name:           "From to end",
			inputText:      "test11\ttest12\ttest13\ttest14\ntest21\ttest22\ttest23\ntest31\ttest32\ttest33\ttest34\ttest35\ntest41 test42\n",
			expectedOutput: "test13\ttest14\t\ntest23\t\ntest33\ttest34\ttest35\t\ntest41 test42\n",
			expectedError:  nil,
			options:        NewOptions(FieldsParams{IsToEnd: true, FromToEnd: 2}, "\t", false),
		}, {
			name:           "From start and to end crossing",
			inputText:      "test11\ttest12\ttest13\ntest21\ttest22\ttest23\ntest31\ttest32\ttest33\n",
			expectedOutput: "test11\ttest12\ttest13\t\ntest21\ttest22\ttest23\t\ntest31\ttest32\ttest33\t\n",
			expectedError:  nil,
			options:        NewOptions(FieldsParams{IsToEnd: true, FromToEnd: 1, IsFromStart: true, FromStartTo: 2}, "\t", false),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.inputText)
			var buffer bytes.Buffer
			err := Cut(reader, &buffer, testCase.options)
			if err != testCase.expectedError {
				t.Errorf("error: got %s, want %s", err, testCase.expectedError)
			}
			if got := buffer.String(); got != testCase.expectedOutput {
				t.Errorf("\ngot \n%swant \n%s", got, testCase.expectedOutput)
			}
		})
	}
}
