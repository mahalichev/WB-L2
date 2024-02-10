package grep

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
)

var ErrNotEnoughArguments error = errors.New("not enough arguments")

type Options struct {
	Filepaths []string
	Pattern   string
	After     int
	Before    int
	Context   int
	Count     bool
	Invert    bool
	LineNum   bool
}

func NewOptions(filepaths []string, pattern string, after, before, context int, count, invert, lineNum bool) Options {
	return Options{
		Filepaths: filepaths,
		Pattern:   pattern,
		After:     after,
		Before:    before,
		Context:   context,
		Count:     count,
		Invert:    invert,
		LineNum:   lineNum,
	}
}

// Получение значений флагов и аргументов
func ParseArguments(arguments []string) (Options, error) {
	fSet := flag.NewFlagSet("grep", flag.ContinueOnError)
	after := fSet.Int("A", 0, "print +N line after match")
	before := fSet.Int("B", 0, "print +N lines before match")
	context := fSet.Int("C", 0, "(A+B) print ±N lines around match")
	count := fSet.Bool("c", false, "number of lines")
	ignoreCase := fSet.Bool("i", false, "ignore case")
	invert := fSet.Bool("v", false, "instead of matching, exclude")
	fixed := fSet.Bool("F", false, "exact match to string, not a pattern")
	lineNum := fSet.Bool("n", false, "print line number")

	if err := fSet.Parse(arguments); err != nil {
		return Options{}, err
	}

	if len(fSet.Args()) < 1 {
		return Options{}, ErrNotEnoughArguments
	}

	pattern := fSet.Arg(0)
	filenames := fSet.Args()[1:]

	// Если не паттерн - экранирование символов
	if *fixed {
		pattern = regexp.QuoteMeta(pattern)
	}
	// Если игнорирование регистра - добавление соответствующего модификатора
	if *ignoreCase {
		pattern = fmt.Sprintf("(?i)%s", pattern)
	}

	return NewOptions(filenames, pattern, *after, *before, *context, *count, *invert, *lineNum), nil
}
