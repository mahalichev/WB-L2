package cut

import (
	"bufio"
	"io"
	"strings"
)

func WriteWithDelimiter(writer *bufio.Writer, content, delimiter string) error {
	if _, err := writer.WriteString(content); err != nil {
		return err
	}
	if _, err := writer.WriteString(delimiter); err != nil {
		return err
	}
	return nil
}

// Реализация утилиты cut
func Cut(in io.Reader, out io.Writer, options Options) error {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	// Получение текста из io.Reader
	text := make([]string, 0)
	for {
		str, err := reader.ReadString('\n')
		str = strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r")
		if err != nil {
			if err != io.EOF {
				return err
			}
			if str != "" {
				text = append(text, str)
			}
			break
		}
		text = append(text, str)
	}

	for _, str := range text {
		// Разбиение строки по разделителю
		splitted := strings.Split(str, options.Delimiter)
		n := len(splitted)

		// Если строка не разделена и разрешено выводить неразделенные строки - вывод строки
		if n == 1 && !options.Separated {
			if _, err := writer.WriteString(str); err != nil {
				return err
			}
			if err := writer.WriteByte('\n'); err != nil {
				return err
			}
			continue
		}

		// Если строка не разделена и запрещено выводить неразделенные строки - пропуск
		if n == 1 && options.Separated {
			continue
		}

		// Если в аргументе -f присутсвовало значение вида "-N" - вывод колонок от 1 до N
		if options.IsFromStart {
			to := min(n-1, options.FromStartTo)
			for i := 0; i <= to; i++ {
				if err := WriteWithDelimiter(writer, splitted[i], options.Delimiter); err != nil {
					return err
				}
			}
		}

		// Вывод отдельных колонок, которые в аргументе -f были записаны либо конкретным числом, либо через интервал
		for _, i := range options.Fields {
			if i >= n {
				continue
			}
			if err := WriteWithDelimiter(writer, splitted[i], options.Delimiter); err != nil {
				return err
			}
		}

		// Если в аргументе -f присутсвовало значение вида "N-" - вывод колонок от N до len(splitted)
		if options.IsToEnd {
			i := options.FromToEnd
			if options.IsFromStart && options.FromToEnd <= options.FromStartTo {
				i = options.FromStartTo + 1
			}
			for ; i <= n-1; i++ {
				if err := WriteWithDelimiter(writer, splitted[i], options.Delimiter); err != nil {
					return err
				}
			}
		}

		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	return nil
}
