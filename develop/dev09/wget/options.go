package wget

import (
	"errors"
	"flag"
	"net/url"
)

var ErrNotEnoughArguments error = errors.New("not enough arguments")
var ErrCantParseURL error = errors.New("can't parse url")

type Options struct {
	BaseURL        *url.URL
	Recursive      bool
	RecursionDepth int
	OtherHosts     bool
}

func NewOptions(baseURL *url.URL, recursive bool, recursionDepth int, otherHosts bool) Options {
	return Options{BaseURL: baseURL, Recursive: recursive, RecursionDepth: recursionDepth, OtherHosts: otherHosts}
}

// Получение значений флагов и аргументов
func ParseArguments(arguments []string) (Options, error) {
	fSet := flag.NewFlagSet("wget", flag.ContinueOnError)
	recursive := fSet.Bool("r", false, "download recursively")
	recursionDepth := fSet.Int("d", 0, "depth of recursion")
	otherHosts := fSet.Bool("h", false, "download from other hosts (in recursive mode)")
	if err := fSet.Parse(arguments); err != nil {
		return Options{}, err
	}

	if len(fSet.Args()) < 1 {
		return Options{}, ErrNotEnoughArguments
	}

	baseURL, err := url.Parse(FormatURL(fSet.Arg(0), ""))
	if err != nil {
		return Options{}, ErrCantParseURL
	}
	return NewOptions(baseURL, *recursive, *recursionDepth, *otherHosts), nil
}
