package sort

import (
	"slices"
	"strings"
)

// Получение значения колонки (вспомогательная функция для OnlyUnique())
func GetColumnValue(text []*StrokeEntry, i, column int, ignoreTrailingBlanks bool) string {
	value := ""
	if len(text[i].Content) > column {
		value = text[i].Content[column]
	}
	if ignoreTrailingBlanks {
		value = strings.TrimRight(value, " ")
	}
	return value
}

// сортировка по колонке
func ColumnSort(text []*StrokeEntry, column int, ignoreTrailingBlanks bool) {
	slices.SortFunc(text, func(a, b *StrokeEntry) int {
		valueA := ""
		if len(a.Content) > column {
			valueA = a.Content[column]
		}
		if ignoreTrailingBlanks {
			valueA = strings.TrimRight(valueA, " ")
		}

		valueB := ""
		if len(b.Content) > column {
			valueB = b.Content[column]
		}
		if ignoreTrailingBlanks {
			valueB = strings.TrimRight(valueB, " ")
		}

		if valueA < valueB {
			return -1
		}
		if valueA > valueB {
			return 1
		}
		return 0
	})
}
