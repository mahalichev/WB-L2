package sort

import (
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type NumericSuffix struct {
	Number int
	Suffix string
}

// Разделение строки на числовое значение и суффикс
func SplitNumericSuffix(str string) NumericSuffix {
	for i, symbol := range str {
		if !unicode.IsDigit(symbol) {
			number, _ := strconv.Atoi(str[:i])
			return NumericSuffix{number, str[i:]}
		}
	}
	number, _ := strconv.Atoi(str)
	return NumericSuffix{number, ""}
}

// Получение числового значения с учетом суффикса по индексу строки и колонке (вспомогательная функция для OnlyUnique())
func GetNumericSuffixesValue(text []*StrokeEntry, i, column int, ignoreTrailingBlanks bool) string {
	if len(text[i].Content) > column {
		word := text[i].Content[column]
		if ignoreTrailingBlanks {
			word = strings.TrimRight(word, " ")
		}
		value := SplitNumericSuffix(word)
		return strconv.Itoa(value.Number) + value.Suffix
	}
	return ""
}

// Сортировка по числовому значению с учетом суффикса
func NumericSuffixesSort(text []*StrokeEntry, column int, ignoreTrailingBlanks bool) {
	slices.SortFunc(text, func(a, b *StrokeEntry) int {
		valueA := NumericSuffix{}
		if len(a.Content) > column {
			value := a.Content[column]
			if ignoreTrailingBlanks {
				value = strings.TrimSuffix(value, " ")
			}
			valueA = SplitNumericSuffix(value)
		}
		valueB := NumericSuffix{}
		if len(b.Content) > column {
			value := b.Content[column]
			if ignoreTrailingBlanks {
				value = strings.TrimSuffix(value, " ")
			}
			valueB = SplitNumericSuffix(value)
		}
		if valueA.Number == valueB.Number {
			if valueA.Suffix < valueB.Suffix {
				return -1
			}
			if valueA.Suffix > valueB.Suffix {
				return 1
			}
			return 0
		}
		return valueA.Number - valueB.Number
	})
}
