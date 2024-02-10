package grep

import (
	"bufio"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Нахождение индексов строк, подходящих под фильтр (с использованием регулярных выражений)
func FindIndexes(text []string, pattern string, invert bool) (map[int]struct{}, error) {
	result := make(map[int]struct{}, 0)
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	for i, str := range text {
		match := compiled.MatchString(str)
		if (!invert && match) || (invert && !match) {
			result[i] = struct{}{}
		}
	}
	return result, nil
}

// Реализация утилиты фильтрации
func GREP(in io.Reader, out io.Writer, options Options) error {
	inReader := bufio.NewReader(in)
	text := make([]string, 0)
	// Получение текста из io.Reader
	for {
		str, err := inReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				str = strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r")
				if str != "" {
					text = append(text, str)
				}
				break
			}
			return err
		}
		text = append(text, strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r"))
	}
	// Поиск индексов строк, подходящих под фильтр
	result, err := FindIndexes(text, options.Pattern, options.Invert)
	if err != nil {
		return err
	}
	// Вывод результата в io.Writer
	if err := GetResult(out, text, result, options.After, options.Before, options.Context, options.Count, options.LineNum); err != nil {
		return err
	}
	return nil
}

func GetResult(out io.Writer, text []string, indices map[int]struct{}, after, before, context int, count, lineNum bool) error {
	outWriter := bufio.NewWriter(out)
	defer outWriter.Flush()

	var err error
	// Если необходимо только количество строк - выводим длину map найденных индексов
	if count {
		if _, err = outWriter.WriteString(strconv.Itoa(len(indices))); err != nil {
			return err
		}
		err = outWriter.WriteByte('\n')
		return err
	}

	// Унифицирование after, before и context
	after = max(after, context)
	before = max(before, context)

	n := len(text)
	// Функция проверки валидности индекса (индекс, не выходящий за пределы text)
	validIndex := func(index int) bool { return 0 <= index && index < n }

	// Множество индексов для вывода
	printSet := make(map[int]struct{})
	for index := range indices {
		if !validIndex(index) {
			continue
		}
		// Добавление индексов в множество с учетом контекста
		for i := index - before; i < index; i++ {
			if !validIndex(i) {
				continue
			}
			printSet[i] = struct{}{}
		}
		printSet[index] = struct{}{}
		for i := index + 1; i <= index+after; i++ {
			if !validIndex(i) {
				continue
			}
			printSet[i] = struct{}{}
		}
	}

	// Сортировка множества индексов для вывода
	printSorted := make([]int, 0, len(printSet))
	for index := range printSet {
		printSorted = append(printSorted, index)
	}
	slices.Sort(printSorted)

	for _, index := range printSorted {
		// Если необходимо добавить номер строки - вывод с номером строки (найденные строки имеют вид N:content, контекст N-content)
		if lineNum {
			outWriter.WriteString(strconv.Itoa(index + 1))
			if _, ok := indices[index]; ok {
				err = outWriter.WriteByte(':')
			} else {
				err = outWriter.WriteByte('-')
			}
			if err != nil {
				return err
			}
		}

		if _, err = outWriter.WriteString(text[index]); err != nil {
			return err
		}

		if err = outWriter.WriteByte('\n'); err != nil {
			return err
		}
	}
	return nil
}
