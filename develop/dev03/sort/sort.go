package sort

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

type StrokeEntry struct {
	Content      []string
	Stroke       string
	InitialIndex int
}

// Удаление специальных символов в конце строки ("\n" или "\r\n")
func trimString(str string) string {
	return strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r")
}

// Получение текста из io.Reader и сохранение его в структуру типа []*StrokeEntry{}
func GetText(in io.Reader) ([]*StrokeEntry, error) {
	reader := bufio.NewReader(in)
	content := []*StrokeEntry{}
	initialIndex := 0
	for {
		str, err := reader.ReadString('\n')
		if err == nil {
			trimmed := trimString(str)
			content = append(content, &StrokeEntry{Content: Split(trimmed, ' '), Stroke: trimmed, InitialIndex: initialIndex})
			initialIndex++
			continue
		}
		if err == io.EOF {
			trimmed := trimString(str)
			if len(trimmed) == 0 {
				return content, nil
			}
			return append(content, &StrokeEntry{Content: Split(trimmed, ' '), Stroke: trimmed, InitialIndex: initialIndex}), nil
		}
		return nil, err
	}
}

// Разделение строки на слова с учётом хвостовых пробелов
func Split(str string, sep rune) []string {
	if len(str) == 0 {
		return nil
	}
	columns := []string{}
	from := 0
	wasSpace := str[0] == ' '
	for i, char := range str {
		if unicode.IsSpace(char) {
			wasSpace = true
			continue
		}
		if wasSpace {
			columns = append(columns, str[from:i-1])
			from = i
			wasSpace = false
		}
	}
	return append(columns, str[from:])
}

// Получение только уникальных строк в структуре типа []*StrokeEntry (уникальная строка - строка, у которой значение, по которому
// производится сортировка, уникально)
func OnlyUnique(text []*StrokeEntry, column int, ignoreTrailingBlanks bool, valueGetter func([]*StrokeEntry, int, int, bool) string) []*StrokeEntry {
	// Реализация типа данных set
	set := make(map[string]struct{})
	for i := 0; i < len(text); i++ {
		// Получение значения строки, по которому будет проводиться сортировка
		value := valueGetter(text, i, column, ignoreTrailingBlanks)
		if _, ok := set[value]; ok {
			n := len(text)
			text[i] = text[n-1]
			text[n-1] = nil
			text = text[:n-1]
			i--
		} else {
			set[value] = struct{}{}
		}
	}
	return text
}

// Сортировка текста
func Sort(in io.Reader, out io.Writer, options Options) error {
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Получение текста из io.Reader
	text, err := GetText(in)
	if err != nil {
		return err
	}
	initLen := len(text)

	switch {
	case options.Numeric:
		if options.Unique {
			text = OnlyUnique(text, options.Column, options.IgnoreTrailingBlanks, GetNumericValue)
		}
		// Сортировка по числовому значению
		NumericSort(text, options.Column, options.IgnoreTrailingBlanks)
	case options.MonthSort:
		if options.Unique {
			text = OnlyUnique(text, options.Column, options.IgnoreTrailingBlanks, GetMonthValue)
		}
		// Сортировка по названию месяца
		MonthSort(text, options.Column, options.IgnoreTrailingBlanks)
	case options.NumericSuffixes:
		if options.Unique {
			text = OnlyUnique(text, options.Column, options.IgnoreTrailingBlanks, GetNumericSuffixesValue)
		}
		// Сортировка по числовому значению с учётом суффиксов
		NumericSuffixesSort(text, options.Column, options.IgnoreTrailingBlanks)
	case options.Column == 0:
		if options.Unique {
			text = OnlyUnique(text, options.Column, options.IgnoreTrailingBlanks, GetDefaultValue)
		}
		// Сортировка по строкам
		DefaultSort(text, options.Column, options.IgnoreTrailingBlanks)
	default:
		if options.Unique {
			text = OnlyUnique(text, options.Column, options.IgnoreTrailingBlanks, GetColumnValue)
		}
		// Сортировка по колонке
		ColumnSort(text, options.Column, options.IgnoreTrailingBlanks)
	}

	// Если в обратном порядке - инверсировать результат
	if options.Reversed {
		n := len(text)
		for i := 0; i < n/2; i++ {
			text[i], text[n-i-1] = text[n-i-1], text[i]
		}
	}

	// Если не проверка на отсортированность io.Reader - вывод результата
	if !options.CheckIfSorted {
		for _, entry := range text {
			if _, err := writer.WriteString(entry.Stroke); err != nil {
				return err
			}
			if err := writer.WriteByte('\n'); err != nil {
				return err
			}
		}
		return nil
	}

	// Проверка на отсортированность io.Reader
	if !IsSorted(text, initLen) {
		if _, err := writer.WriteString("not sorted\n"); err != nil {
			return err
		}
	}
	return nil
}

func IsSorted(text []*StrokeEntry, initLen int) bool {
	// Если длина исходного текста и результата не идентична - не отсортирован
	if len(text) != initLen {
		return false
	}
	// Проверка, что каждая строка осталась на своём месте
	for i, entry := range text {
		if i != entry.InitialIndex {
			return false
		}
	}
	return true
}
