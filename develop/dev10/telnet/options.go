package telnet

import (
	"errors"
	"flag"
	"time"
)

var ErrNotEnoughArguments error = errors.New("not enough arguments")

type Options struct {
	Host    string
	Port    string
	Timeout time.Duration
}

func NewOptions(host, port string, timeout time.Duration) Options {
	return Options{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
}

// Получение значений флагов и аргументов
func ParseArguments(arguments []string) (Options, error) {
	fSet := flag.NewFlagSet("telnet", flag.ContinueOnError)
	timeout := fSet.Duration("timeout", 10*time.Second, "connection timeout")
	if err := fSet.Parse(arguments); err != nil {
		return Options{}, err
	}

	if len(fSet.Args()) < 2 {
		return Options{}, ErrNotEnoughArguments
	}

	return NewOptions(fSet.Arg(0), fSet.Arg(1), *timeout), nil
}
