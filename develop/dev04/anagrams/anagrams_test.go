package anagrams

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		expected map[string][]string
	}{
		{
			name:     "Standard test from the condition",
			words:    []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			expected: map[string][]string{"пятак": {"пятак", "пятка", "тяпка"}, "листок": {"листок", "слиток", "столик"}},
		}, {
			name:     "Anagrams slice of length 1",
			words:    []string{"пятак", "листок", "слиток", "столик"},
			expected: map[string][]string{"листок": {"листок", "слиток", "столик"}},
		}, {
			name:     "Empty words slice",
			words:    []string{},
			expected: map[string][]string{},
		}, {
			name:     "Repeats in words slice + case sensitivity",
			words:    []string{"пятак", "пяТАК", "пятка", "тяпка", "ЛИСТОК", "СТОЛик"},
			expected: map[string][]string{"пятак": {"пятак", "пятка", "тяпка"}, "листок": {"листок", "столик"}},
		}, {
			name:     "Words order",
			words:    []string{"пятка", "столик", "тяпка", "пятак", "листок", "слиток"},
			expected: map[string][]string{"пятка": {"пятак", "пятка", "тяпка"}, "столик": {"листок", "слиток", "столик"}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := AnagramsMap(&testCase.words)
			if !reflect.DeepEqual(*got, testCase.expected) {
				t.Errorf("got %v, want %v", *got, testCase.expected)
			}
		})
	}
}

func TestSortLetters(t *testing.T) {
	testCases := []struct {
		name     string
		word     string
		expected string
	}{
		{
			name:     "Common sort",
			word:     "дгвба",
			expected: "абвгд",
		}, {
			name:     "Sorted string",
			word:     "абвгд",
			expected: "абвгд",
		}, {
			name:     "Empty string",
			word:     "",
			expected: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := SortLetters(testCase.word)
			if got != testCase.expected {
				t.Errorf("got %v, want %v", got, testCase.expected)
			}
		})
	}
}

func TestClearAnagramsMap(t *testing.T) {
	testCases := []struct {
		name        string
		anagramsMap map[string][]string
		expected    map[string][]string
	}{
		{
			name:        "No deletions",
			anagramsMap: map[string][]string{"пятак": {"пятак", "пятка", "тяпка"}, "листок": {"листок", "слиток", "столик"}},
			expected:    map[string][]string{"пятак": {"пятак", "пятка", "тяпка"}, "листок": {"листок", "слиток", "столик"}},
		}, {
			name:        "Slice with length 1 deletion",
			anagramsMap: map[string][]string{"пятак": {"пятак"}, "листок": {"листок", "слиток", "столик"}},
			expected:    map[string][]string{"листок": {"листок", "слиток", "столик"}},
		}, {
			name:        "Empty map",
			anagramsMap: map[string][]string{},
			expected:    map[string][]string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ClearAnagramsMap(testCase.anagramsMap)
			if !reflect.DeepEqual(testCase.anagramsMap, testCase.expected) {
				t.Errorf("got %v, want %v", testCase.anagramsMap, testCase.expected)
			}
		})
	}
}

func TestSortAnagramsMap(t *testing.T) {
	testCases := []struct {
		name        string
		anagramsMap map[string][]string
		expected    map[string][]string
	}{
		{
			name:        "Default map",
			anagramsMap: map[string][]string{"пятак": {"пятак", "тяпка", "пятка"}, "листок": {"слиток", "столик", "листок"}},
			expected:    map[string][]string{"пятак": {"пятак", "пятка", "тяпка"}, "листок": {"листок", "слиток", "столик"}},
		}, {
			name:        "Empty map",
			anagramsMap: map[string][]string{},
			expected:    map[string][]string{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			SortAnagramsMap(testCase.anagramsMap)
			if !reflect.DeepEqual(testCase.anagramsMap, testCase.expected) {
				t.Errorf("got %v, want %v", testCase.anagramsMap, testCase.expected)
			}
		})
	}
}
