package shell

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

var ErrNotEnoughArguments error = errors.New("not enough arguments")
var ErrNotEnoughCommands error = errors.New("not enough commands in pipe")

// Реализация shell
func Shell(in io.Reader, out io.Writer) error {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	for {
		// Получение очередной строки
		stroke, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		// Удаление символов \r\n в конце строки
		stroke = strings.TrimSuffix(strings.TrimSuffix(stroke, "\n"), "\r")

		// Если введено "\quit" - завершение работы функции
		if stroke == "\\quit" {
			return nil
		}

		var result string
		// Если есть пайп - запуск конвеера
		if strings.Contains(stroke, " | ") {
			result, err = Pipe(stroke)
			// Иначе запуск команды
		} else {
			result, err = Command(stroke, in)
		}
		// Если возникла ошибка - вывод ошибки
		if err != nil {
			if _, err := writer.WriteString("error: " + err.Error()); err != nil {
				return err
			}
			if err := writer.WriteByte('\n'); err != nil {
				return err
			}
		}
		// Вывод результата работы команды
		if result != "" {
			if _, err := writer.WriteString(result); err != nil {
				return err
			}
			if err := writer.WriteByte('\n'); err != nil {
				return err
			}
		}
		writer.Flush()
	}
}

// Выполнение команды
func Command(stroke string, in io.Reader) (string, error) {
	arguments := strings.Fields(stroke)
	// Если была введена пустая строка - выход из функции
	if len(arguments) == 0 {
		return "", nil
	}

	var result string
	var err error

	// Выполнение команды. Для выполнения команд cd/pwd/ps/kill/echo через exec необходимо в начале строки указать "unix"
	// Например: unix echo hello world
	switch arguments[0] {
	case "cd":
		directory := ""
		if len(arguments) > 1 {
			directory = arguments[1]
		}
		err = Cd(directory)
	case "pwd":
		result, err = Pwd()
	case "ps":
		result, err = Ps()
	case "kill":
		err = Kill(arguments[1:]...)
	case "echo":
		result = Echo(arguments[1:]...)
	case "unix":
		if len(arguments) < 2 {
			return "", ErrNotEnoughArguments
		}
		result, err = Exec(arguments[1], in, arguments[2:]...)
	default:
		result, err = Exec(arguments[0], in, arguments[1:]...)
	}
	return strings.TrimSuffix(result, "\n"), err
}

// Реализация конвеера на пайпах
func Pipe(stroke string) (string, error) {
	commands := strings.Split(stroke, " | ")
	// Если введено менее двух команд
	if strings.HasSuffix(strings.TrimSpace(stroke), "|") || len(commands) < 2 {
		return "", ErrNotEnoughCommands
	}
	// Использование буффера для сохранения вывода одной команды и передачи вводом в другую
	var inBuffer bytes.Buffer
	for _, command := range commands {
		arguments := strings.Fields(command)
		// Если команда не найдена
		if len(arguments) == 0 {
			return "", ErrNotEnoughCommands
		}
		result, err := Exec(arguments[0], &inBuffer, arguments[1:]...)
		if err != nil {
			return "", err
		}
		if _, err := inBuffer.WriteString(result); err != nil {
			return "", err
		}
	}
	return strings.TrimSuffix(inBuffer.String(), "\n"), nil
}
