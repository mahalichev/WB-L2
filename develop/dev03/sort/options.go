package sort

import (
	"errors"
	"flag"
)

var ErrNonPositiveColumn error = errors.New("column must be a positive number")
var ErrNotEnoughArguments error = errors.New("not enough arguments")

type Options struct {
	Filepath             string
	Column               int
	Numeric              bool
	MonthSort            bool
	NumericSuffixes      bool
	Reversed             bool
	Unique               bool
	IgnoreTrailingBlanks bool
	CheckIfSorted        bool
}

func NewOptions(filepath string, column int, numeric, monthSort, numericSuffixes, reversed, unique, ignoreTrailingBlanks, checkIfSorted bool) Options {
	return Options{
		Filepath:             filepath,
		Column:               column,
		Numeric:              numeric,
		MonthSort:            monthSort,
		NumericSuffixes:      numericSuffixes,
		Reversed:             reversed,
		Unique:               unique,
		IgnoreTrailingBlanks: ignoreTrailingBlanks,
		CheckIfSorted:        checkIfSorted,
	}
}

// Получение значений флагов и аргументов
func ParseArguments(arguments []string) (Options, error) {
	fSet := flag.NewFlagSet("sort", flag.ContinueOnError)
	column := fSet.Int("k", 1, "column for sorting")
	numeric := fSet.Bool("n", false, "sort by numeric value")
	reversed := fSet.Bool("r", false, "sort in reverse order")
	unique := fSet.Bool("u", false, "only unique strings")
	monthSort := fSet.Bool("M", false, "sort by month name")
	ignoreTrailingBlanks := fSet.Bool("b", false, "ignore trailing spaces")
	checkIfSorted := fSet.Bool("c", false, "check if data is sorted")
	numericSuffixes := fSet.Bool("h", false, "sort by numeric value taking into account suffixes")
	if err := fSet.Parse(arguments); err != nil {
		return Options{}, err
	}

	if len(fSet.Args()) < 1 {
		return Options{}, ErrNotEnoughArguments
	}
	filepath := fSet.Arg(0)

	if *column < 1 {
		return Options{}, ErrNonPositiveColumn
	}
	return NewOptions(filepath, *column-1, *numeric, *monthSort, *numericSuffixes, *reversed, *unique, *ignoreTrailingBlanks, *checkIfSorted), nil
}
