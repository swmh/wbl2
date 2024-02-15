package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

type Match struct {
	start, len int
}

type Entry struct {
	matches []Match
	b       []byte
	line    int
}

func Grep(matcher Matcher, printer Printer, r io.Reader) {
	scanner := bufio.NewScanner(r)

	for i := 0; scanner.Scan(); i++ {
		matches := matcher.Match(scanner.Bytes())
		en := Entry{
			b:       scanner.Bytes(),
			line:    i,
			matches: matches,
		}

		printer.Print(en)
	}

	printer.Result()
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	m, err := cfg.GetMatcher()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(cfg.Inputs) == 0 {
		pr, err := cfg.GetPrinter(os.Stdout)
		if err != nil {
			fmt.Println(err)
			return
		}
		Grep(m, pr, os.Stdin)
		return
	}

	for _, v := range cfg.Inputs {
		f, err := os.Open(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		pr, err := cfg.GetPrinter(os.Stdout)
		if err != nil {
			fmt.Println(err)
			return
		}

		Grep(m, pr, f)
	}
}
