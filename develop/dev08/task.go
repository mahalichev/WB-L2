package main

import (
	"dev08/shell"
	"fmt"
	"os"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	if err := shell.Shell(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
