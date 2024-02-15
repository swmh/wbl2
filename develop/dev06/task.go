package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*
=== Утилита cut ===

# Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
type Config struct {
	Ranges []Range

	Delimiter string

	OnlySeparated bool
}

func MustParse() Config {
	c := Config{}

	var fields string

	flag.StringVar(&fields, "f", "", "")
	flag.StringVar(&c.Delimiter, "d", "\t", "")
	flag.BoolVar(&c.OnlySeparated, "s", false, "")

	flag.Parse()


	if fields == "" {
		log.Fatalln("cut: require an argument \"-f\"")
	}

	if utf8.RuneCountInString(c.Delimiter) != 1 {
		log.Fatalln("cut: delimiter must be single character")
	}

	parts := strings.Split(fields, ",")

	for _, part := range parts {
		if strings.Contains(part, "-") {
			ran, err := ParseRange(part)
			if err != nil {
				log.Fatalln(err)
			}

			c.Ranges = append(c.Ranges, ran)
			continue
		}

		n, err := strconv.ParseUint(part, 10, 0)
		if err != nil {
			log.Fatalln(errInvalidFieldValue)
		}

		c.Ranges = append(c.Ranges, Range{
			start: n - 1,
			end:   n - 1,
		})

	}

	return c
}

var (
	errInvalidFieldValue = errors.New("invalid field value")
	errInvalidRange      = errors.New("invalid range")
)

type Range struct {
	start uint64
	end   uint64
}

func (r Range) InRange(n int) bool {
	if n < 0 {
		return false
	}

	v := uint64(n)
	return v >= r.start && v <= r.end
}

func ParseRange(s string) (Range, error) {
	part := strings.Split(s, "-")

	if len(part) != 2 {
		return Range{}, errInvalidRange
	}

	if part[0] == "" && part[1] == "" {
		return Range{}, errInvalidRange
	}

	ran := Range{
		start: 0,
		end:   ^uint64(0),
	}

	if part[0] != "" {
		v, err := strconv.ParseUint(part[0], 10, 0)
		if v == 0 {
			return Range{}, errInvalidRange
		}

		ran.start = v - 1

		if err != nil {
			return Range{}, errInvalidRange
		}
	}

	if part[1] != "" {
		v, err := strconv.ParseUint(part[1], 10, 0)
		if v == 0 {
			return Range{}, errInvalidRange
		}

		ran.end = v - 1

		if err != nil {
			return Range{}, errInvalidRange
		}

	}

	return ran, nil
}

func main() {
	cfg := MustParse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), cfg.Delimiter)
		if len(splitted) == 1 {
			if cfg.OnlySeparated {
				continue
			}

			fmt.Println(splitted)
		}

		var fields []string

		for i, field := range splitted {
			for _, r := range cfg.Ranges {
				if r.InRange(i) {
					fields = append(fields, field)
					break
				}
			}
		}

		fmt.Println(strings.Join(fields, cfg.Delimiter))

	}
}
