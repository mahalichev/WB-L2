package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	gops "github.com/mitchellh/go-ps"
)

// Реализация команды cd
func Cd(directory string) error {
	// Если директория не указана - перемещение на домашнюю директорию пользователя
	if directory == "" {
		directory, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return os.Chdir(directory)
	}
	// Изменение директории на указанную
	return os.Chdir(directory)
}

// Реализация команды echo
func Echo(arguments ...string) string {
	return strings.Join(arguments, " ")
}

// Реализация команды pwd
func Pwd() (string, error) {
	// Получение текущей директории
	return os.Getwd()
}

// Реализация команды ps
func Ps() (string, error) {
	// Получение информации о процессах
	processes, err := gops.Processes()
	if err != nil {
		return "", err
	}
	var builder strings.Builder
	// Запись информации в строку (ID процесса, ID родительского процесса и исполняемый файл)
	for _, process := range processes {
		builder.WriteString(fmt.Sprintf("%d %d %s\n", process.Pid(), process.PPid(), process.Executable()))
	}
	return builder.String(), nil
}

// Реализация команды kill
func Kill(pids ...string) error {
	for _, pid := range pids {
		pid, err := strconv.Atoi(pid)
		if err != nil {
			return err
		}
		// Поиск процесса с заданным ID
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		// Завершение процесса
		if err := process.Kill(); err != nil {
			return err
		}
	}
	return nil
}

// Выполнение fork/exec команды
func Exec(command string, in io.Reader, arguments ...string) (string, error) {
	cmd := exec.Command(command, arguments...)
	var resBuf, errBuf bytes.Buffer
	// Установка буферов в качестве потоков ввода/вывода информации
	cmd.Stdin = in
	cmd.Stderr = &errBuf
	cmd.Stdout = &resBuf

	// Запуск и ожидание завершения работы команды
	if err := cmd.Run(); err != nil {
		return "", err
	}

	if errBuf.Len() == 0 {
		return resBuf.String(), nil
	}
	return resBuf.String(), errors.New(errBuf.String())
}
