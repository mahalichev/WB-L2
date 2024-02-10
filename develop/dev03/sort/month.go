package sort

import (
	"slices"
	"strconv"
	"strings"
)

// Порядок месяцев
var MonthsOrder map[string]int = map[string]int{
	"Jan":       11,
	"January":   12,
	"Feb":       21,
	"February":  22,
	"Mar":       31,
	"March":     32,
	"Apr":       41,
	"April":     42,
	"May":       51,
	"Jun":       61,
	"June":      62,
	"Jul":       71,
	"July":      72,
	"Aug":       81,
	"August":    82,
	"Sep":       91,
	"Sept":      92,
	"September": 93,
	"Oct":       101,
	"October":   102,
	"Nov":       111,
	"November":  112,
	"Dec":       121,
	"December":  122,
}

// Получение месяца по индексу строки и колонке (вспомогательная функция для OnlyUnique())
func GetMonthValue(text []*StrokeEntry, i, column int, ignoreTrailingBlanks bool) string {
	if len(text[i].Content) > column {
		key := text[i].Content[column]
		if ignoreTrailingBlanks {
			key = strings.TrimRight(key, " ")
		}
		if val, ok := MonthsOrder[key]; ok {
			return strconv.Itoa(val / 10)
		}
	}
	return ""
}

// Сортировка по месяцу
func MonthSort(text []*StrokeEntry, column int, ignoreTrailingBlanks bool) {
	slices.SortFunc(text, func(a, b *StrokeEntry) int {
		valueA := 0
		if len(a.Content) > column {
			key := a.Content[column]
			if ignoreTrailingBlanks {
				key = strings.TrimSuffix(key, " ")
			}
			if val, ok := MonthsOrder[key]; ok {
				valueA = val
			}
		}
		valueB := 0
		if len(b.Content) > column {
			key := b.Content[column]
			if ignoreTrailingBlanks {
				key = strings.TrimSuffix(key, " ")
			}
			if val, ok := MonthsOrder[key]; ok {
				valueB = val
			}
		}
		return valueA - valueB
	})
}
