package ntp

import (
	"time"

	"github.com/beevik/ntp"
)

var DefaultNTPAddress = "0.beevik-ntp.pool.ntp.org"

func GetNTPTime(address string) (time.Time, error) {
	if address == "" {
		address = DefaultNTPAddress
	}
	// Получение точного времени с использованием библиотеки github.com/beevik/ntp
	currentTime, err := ntp.Time(address)
	if err != nil {
		return time.Time{}, err
	}
	return currentTime, err
}
