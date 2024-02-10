package main

import (
	"fmt"
	"os"

	"dev06/cut"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	options, err := cut.ParseArguments(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if err := cut.Cut(os.Stdin, os.Stdout, options); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
