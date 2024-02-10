package sort

import (
	"slices"
	"strconv"
	"strings"
)

// Получение числа по индексу строки и колонке (вспомогательная функция для OnlyUnique())
func GetNumericValue(text []*StrokeEntry, i, column int, ignoreTrailingBlanks bool) string {
	if len(text[i].Content) > column {
		word := text[i].Content[column]
		if ignoreTrailingBlanks {
			word = strings.TrimRight(word, " ")
		}
		if value, err := strconv.Atoi(word); err == nil {
			return strconv.Itoa(value)
		}
	}
	return ""
}

// Сортировка по числу
func NumericSort(text []*StrokeEntry, column int, ignoreTrailingBlanks bool) {
	slices.SortFunc(text, func(a, b *StrokeEntry) int {
		valueA := 0
		if len(a.Content) > column {
			val := a.Content[column]
			if ignoreTrailingBlanks {
				val = strings.TrimSuffix(val, " ")
			}
			valueA, _ = strconv.Atoi(val)
		}
		valueB := 0
		if len(b.Content) > column {
			val := b.Content[column]
			if ignoreTrailingBlanks {
				val = strings.TrimSuffix(val, " ")
			}
			valueB, _ = strconv.Atoi(val)
		}
		return valueA - valueB
	})
}
