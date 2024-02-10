package telnet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
)

// Реализация утилиты Telnet
func Telnet(in io.Reader, out io.Writer, options Options) error {
	// Подключение к адресу по протоколу tcp
	connection, err := net.DialTimeout("tcp", net.JoinHostPort(options.Host, options.Port), options.Timeout)
	if err != nil {
		return err
	}
	defer connection.Close()

	// Перехват сигнала завершения работы программы
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		inReader := bufio.NewReader(in)
		connectionToOut := bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(out))
		for {
			// Получение очередной строки из io.Reader
			data, err := inReader.ReadBytes('\n')
			if err != nil {
				// Если завершение ввода - выход из горутины
				if err == io.EOF {
					return
				}
				// Вывод возникшей ошибки
				fmt.Fprintln(os.Stderr, err)
			}
			// Запись данных в соединение
			if _, err := connection.Write(data); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			// Чтение данных из соединения
			if data, err = connectionToOut.Reader.ReadBytes('\n'); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

			// Запись данных данных в io.Writer
			if _, err := connectionToOut.Writer.Write(data); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			// Освобождение буфера
			if err := connectionToOut.Writer.Flush(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}()

	// Ожидание разрыва подключения/сигнала о завершении работы программы
	select {
	case <-done:
	case <-c:
		signal.Stop(c)
		close(c)
	}
	return nil
}
