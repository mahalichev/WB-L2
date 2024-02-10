package main

import (
	"fmt"
	"os"

	"dev05/grep"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	options, err := grep.ParseArguments(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if len(options.Filepaths) == 0 {
		if err := grep.GREP(os.Stdin, os.Stdout, options); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	for _, filepath := range options.Filepaths {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		defer file.Close()
		fmt.Println(filepath)
		if err := grep.GREP(file, os.Stdout, options); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
