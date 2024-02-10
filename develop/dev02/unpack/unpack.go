package unpack

import (
	"errors"
	"strings"
	"unicode"
)

var ErrIncorrectPackedString = errors.New("packed string is incorrect")

// Продублировать letter count раз
func RepeatLetter(letter string, count int) string {
	if count < 1 {
		count = 1
	}
	return strings.Repeat(letter, count)
}

func UnpackString(packedString string) (string, error) {
	// Builder для эффективной конкатенации строк
	builder := &strings.Builder{}
	letter := ""
	isEscaping := false
	count := 0
	// Прохождение по каждому символу строки
	for _, symbol := range packedString {
		// Если символ - цифра
		if unicode.IsDigit(symbol) {
			// Если необходим escape
			if isEscaping {
				// Устанавливается символ как тот, который будет продублирован
				letter = string(symbol)
				count = 0
				isEscaping = false
				continue
			}
			// Если не было символа для дублирования, возвращается ошибка о некорректной строке
			if letter == "" {
				return "", ErrIncorrectPackedString
			}
			// Если текущий символ 0 и является первой цифрой числа, возвращается ошибка о некорректной строке
			if symbol == '0' && count == 0 {
				return "", ErrIncorrectPackedString
			}
			// Рассчет текущего count
			count = count*10 + int(symbol-'0')
			continue
		}

		// Если символ - \
		if symbol == '\\' {
			// Если необходим escape
			if isEscaping {
				// Устанавливается символ как тот, который будет продублирован
				letter = "\\"
				count = 0
				isEscaping = false
				continue
			}

			// Если до этого символа встречались символы, которые необходимо продублировать
			if letter != "" {
				// Дублируются ранее прочитанные символы
				builder.WriteString(RepeatLetter(letter, count))
				letter = ""
				count = 0
			}

			// Установка флага о необходимости escape
			isEscaping = true
			continue
		}
		// Если символ - любой другой символ и необходимо произвести escape, возвращается ошибка о некорректной строке
		if isEscaping {
			return "", ErrIncorrectPackedString
		}
		// Если до этого был символ, который не был продублирован - дублирование символа
		if letter != "" {
			builder.WriteString(RepeatLetter(letter, count))
		}

		// Установка символа как необходимого к дублированию
		letter = string(symbol)
		count = 0
	}

	// Если остался непродублированный символ - дублирование символа
	if letter != "" {
		builder.WriteString(RepeatLetter(letter, count))
		return builder.String(), nil
	}

	// Если был символ \ без escape символа, возвращается ошибка о некорректной строке
	if isEscaping {
		return "", ErrIncorrectPackedString
	}
	return builder.String(), nil
}
