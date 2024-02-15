package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Sort struct {
	Column      int
	ByNum       bool
	Desc        bool
	PrintUnique bool

	ByMonth       bool
	IgnoreTail    bool
	CheckIsSorted bool
	WithSuffix    bool

	SortFunc func(a string, b string) int
}

func NewSort() Sort {
	s := Sort{}

	flag.IntVar(&s.Column, "k", 0, "")
	flag.BoolVar(&s.ByNum, "n", false, "")
	flag.BoolVar(&s.Desc, "r", false, "")
	flag.BoolVar(&s.PrintUnique, "u", false, "")

	flag.BoolVar(&s.ByMonth, "M", false, "")
	flag.BoolVar(&s.IgnoreTail, "b", false, "")
	flag.BoolVar(&s.CheckIsSorted, "s", false, "")
	flag.BoolVar(&s.WithSuffix, "h", false, "")

	flag.Parse()

	var sortFlags int

	if s.ByMonth {
		sortFlags++
		s.SortFunc = SortByMonth
	}

	if s.ByNum {
		sortFlags++
		s.SortFunc = SortByNum
	}

	if s.WithSuffix {
		sortFlags++
		s.SortFunc = SortWithSuffix
	}

	if sortFlags > 1 {
		fmt.Println("only one sort flag")
		os.Exit(0)
	}

	if s.SortFunc == nil {
		s.SortFunc = cmp.Compare[string]
	}

	return s
}

const whitespace = " \t\n\v\f\r"

func (s *Sort) IsSorted(x []string, sort func(a, b string) int) bool {
	return slices.IsSortedFunc(x, func(a, b string) int {
		if s.Desc {
			a, b = b, a
		}

		a = GetColumn(a, s.Column)
		b = GetColumn(b, s.Column)

		if s.IgnoreTail {
			a = strings.TrimRight(a, whitespace)
			b = strings.TrimRight(b, whitespace)
		}

		return sort(a, b)
	})
}

func (s *Sort) Sort(x []string, sort func(a, b string) int) {
	slices.SortFunc(x, func(a, b string) int {
		if s.Desc {
			a, b = b, a
		}

		a = GetColumn(a, s.Column)
		b = GetColumn(b, s.Column)

		if s.IgnoreTail {
			a = strings.TrimRight(a, whitespace)
			b = strings.TrimRight(b, whitespace)
		}

		return sort(a, b)
	})
}

func UniqueValues(x []string) []string {
	resMap := make(map[string]bool)
	var result []string

	for _, v := range x {
		if resMap[v] {
			continue
		}

		resMap[v] = true
		result = append(result, v)
	}

	return result
}

func GetColumn(x string, n int) string {
	v := strings.Split(x, " ")
	if len(v)-1 < n {
		return ""
	}

	return v[n]
}

func SortByNum(a, b string) int {
	vA, aerr := strconv.ParseFloat(a, 64)
	vB, berr := strconv.ParseFloat(b, 64)

	if aerr != nil && berr != nil {
		return cmp.Compare(a, b)
	}

	if aerr != nil {
		return -1
	}

	if berr != nil {
		return 1
	}

	return cmp.Compare(vA, vB)
}

var suffixes = map[string]int{
	"K":  1_000,
	"KB": 1024,

	"M":  1_000_000,
	"MB": 1_048_576,

	"G":  1_000_000_000,
	"GB": 1_073_741_824,
}

func GetNumSuffix(x string) (string, int) {
	for k, v := range suffixes {
		before, ok := strings.CutSuffix(x, k)
		if ok {
			return before, v
		}
	}

	return "", 0
}

func SortWithSuffix(a, b string) int {
	numA, s1 := GetNumSuffix(a)
	if s1 == 0 {
		s1 = 1
		numA = a
	}

	numB, s2 := GetNumSuffix(b)
	if s2 == 0 {
		s2 = 1
		numB = b
	}

	vA, aerr := strconv.ParseFloat(numA, 64)
	vB, berr := strconv.ParseFloat(numB, 64)

	if aerr != nil && berr != nil {
		return cmp.Compare(a, b)
	}

	if aerr != nil {
		return -1
	}

	if berr != nil {
		return 1
	}

	return cmp.Compare(vA*float64(s1), vB*float64(s2))
}

var months = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,
}

func SortByMonth(a, b string) int {
	vA := months[strings.ToLower(a)]
	vB := months[strings.ToLower(b)]

	return cmp.Compare(vA, vB)
}

func main() {
	cfg := NewSort()

	var arr []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		arr = append(arr, scanner.Text())
	}

	if cfg.CheckIsSorted {
		fmt.Println(cfg.IsSorted(arr, cfg.SortFunc))
		return
	}

	cfg.Sort(arr, cfg.SortFunc)

	if cfg.PrintUnique {
		arr = UniqueValues(arr)
	}

	for _, v := range arr {
		fmt.Println(v)
	}
}
