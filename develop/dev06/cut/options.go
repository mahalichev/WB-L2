package cut

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

var ErrParsingError error = errors.New("an error occurred while parsing flags")
var ErrWrongFieldEntry error = errors.New("-f must me string with positive numbers like `1,2,3`, `1-3`, `-3`, `5-`, or combined `1,3-4,6,10-`")

type FieldsParams struct {
	Fields      []int
	FromStartTo int
	FromToEnd   int
	IsFromStart bool
	IsToEnd     bool
}

// Приведение значение флага -f к структуре типа FieldsParams
func GetFieldsFromString(fieldsStr string) (FieldsParams, error) {
	var err error
	fieldsParams := FieldsParams{}
	intervals := strings.Split(fieldsStr, ",")
	fieldsParams.Fields = make([]int, 0, len(intervals))
	fieldsSet := make(map[int]struct{})

	// Обработка значения вида "-N"
	if strings.HasPrefix(intervals[0], "-") {
		fieldsParams.FromStartTo, err = strconv.Atoi(strings.TrimPrefix(intervals[0], "-"))
		if err != nil || fieldsParams.FromStartTo < 1 {
			return FieldsParams{}, ErrWrongFieldEntry
		}
		fieldsParams.FromStartTo--
		fieldsParams.IsFromStart = true
		intervals = intervals[1:]
	}

	// Обработка значения вида "N-"
	if n := len(intervals) - 1; n >= 0 && strings.HasSuffix(intervals[n], "-") {
		fieldsParams.FromToEnd, err = strconv.Atoi(strings.TrimSuffix(intervals[n], "-"))
		if err != nil || fieldsParams.FromToEnd < 1 {
			return FieldsParams{}, ErrWrongFieldEntry
		}
		fieldsParams.FromToEnd--
		fieldsParams.IsToEnd = true
		intervals = intervals[:n]
	}

	// Обработка остальных значений
	for _, interval := range intervals {
		splittedInterval := strings.Split(interval, "-")
		// Если очередное значение - интервал
		if len(splittedInterval) == 2 {
			start, err := strconv.Atoi(splittedInterval[0])
			if err != nil || start < 1 {
				return FieldsParams{}, ErrWrongFieldEntry
			}

			end, err := strconv.Atoi(splittedInterval[1])
			if err != nil || start < 1 {
				return FieldsParams{}, ErrWrongFieldEntry
			}

			for ; start <= end; start++ {
				value := start - 1
				if _, ok := fieldsSet[value]; ok {
					continue
				}
				if fieldsParams.IsFromStart && value <= fieldsParams.FromStartTo {
					continue
				}
				if fieldsParams.IsToEnd && value >= fieldsParams.FromToEnd {
					continue
				}
				fieldsParams.Fields = append(fieldsParams.Fields, value)
				fieldsSet[value] = struct{}{}
			}
			// Если очередное значение - число
		} else if len(splittedInterval) == 1 {
			number, err := strconv.Atoi(splittedInterval[0])
			if err != nil || number < 1 {
				return FieldsParams{}, ErrWrongFieldEntry
			}
			value := number - 1
			if _, ok := fieldsSet[value]; ok {
				continue
			}
			if fieldsParams.IsFromStart && value <= fieldsParams.FromStartTo {
				continue
			}
			if fieldsParams.IsToEnd && value >= fieldsParams.FromToEnd {
				continue
			}
			fieldsParams.Fields = append(fieldsParams.Fields, value)
			fieldsSet[value] = struct{}{}
		} else {
			return FieldsParams{}, ErrWrongFieldEntry
		}
	}
	return fieldsParams, nil
}

type Options struct {
	FieldsParams
	Delimiter string
	Separated bool
}

func NewOptions(fieldsParams FieldsParams, delimiter string, separated bool) Options {
	return Options{
		FieldsParams: fieldsParams,
		Delimiter:    delimiter,
		Separated:    separated,
	}
}

// Получение значений флагов и аргументов
func ParseArguments(arguments []string) (Options, error) {
	fSet := flag.NewFlagSet("cut", flag.ContinueOnError)
	fieldsStr := fSet.String("f", "", "fields for output")
	delimiter := fSet.String("d", "\t", "delimiter for input")
	separated := fSet.Bool("s", false, "output only separated input")
	if fSet.Parse(arguments) != nil {
		return Options{}, ErrParsingError
	}
	// Приведение флага -f к структуре типа FieldsParams
	fields, err := GetFieldsFromString(*fieldsStr)
	if err != nil {
		return Options{}, err
	}
	return NewOptions(fields, *delimiter, *separated), nil
}
