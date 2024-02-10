package sort

import (
	"slices"
	"strings"
)

// Получение строки по индексу (вспомогательная функция для OnlyUnique())
func GetDefaultValue(text []*StrokeEntry, i, column int, ignoreTrailingBlanks bool) string {
	if ignoreTrailingBlanks {
		return strings.TrimRight(text[i].Stroke, " ")
	}
	return text[i].Stroke
}

// Сортировка по строкам
func DefaultSort(text []*StrokeEntry, column int, ignoreTrailingBlanks bool) {
	slices.SortFunc(text, func(a, b *StrokeEntry) int {
		if a.Stroke < b.Stroke {
			return -1
		}
		if a.Stroke > b.Stroke {
			return 1
		}
		return 0
	})
}
