package main

import (
	"fmt"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

// Создание нового сигнала
func newSignal(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// Реализация функции or
func or(channels ...<-chan interface{}) <-chan interface{} {
	result := make(chan interface{})

	// Создание WaitGroup для ожидания выполнения работы всех горутин
	wg := sync.WaitGroup{}
	// Увеличение счётчика горутин на количество каналов, посылающих сигнал
	wg.Add(len(channels))
	for _, channel := range channels {
		// Запуск горутины для каждого канала, которая будет отправлять сигнал в result при сигнале соответствующего канала
		go func(ch <-chan interface{}) {
			// Уменьшение счётчика горутин
			defer wg.Done()
			result <- <-ch
		}(channel)
	}

	// Запуск горутины, которая закроет канал после получения всех сигналов
	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func main() {
	start := time.Now()
	<-or(
		newSignal(2*time.Hour),
		newSignal(5*time.Minute),
		newSignal(1*time.Second),
		newSignal(1*time.Hour),
		newSignal(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
