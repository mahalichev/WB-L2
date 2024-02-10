package main

import (
	"dev03/sort"
	"fmt"
	"os"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	options, err := sort.ParseArguments(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	file, err := os.Open(options.Filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer file.Close()
	if err := sort.Sort(file, os.Stdout, options); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
