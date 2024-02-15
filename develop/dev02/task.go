package main

import (
	"fmt"
	"strconv"
)

const escapeRune = '\\'

func Escape(r rune) (rune, error) {
	var err error

	_, err = strconv.Atoi(string(r))
	if err == nil {
		return r, nil
	}

	v, _, _, err := strconv.UnquoteChar(string(escapeRune)+string(r), 0)
	if err != nil {
		return 0, err
	}

	return v, err
}

func Unpack(s string) (string, error) {
	var esc bool
	var err error
	var st []rune

	for _, v := range s {
		if !esc && v == escapeRune {
			esc = true
			continue
		}

		if esc {
			v, err = Escape(v)
			if err != nil {
				return "", fmt.Errorf("cannot escape rune: %s", string(v))
			}

			st = append(st, v)
			esc = false

			continue
		}

		var n int
		n, err = strconv.Atoi(string(v))
		if err != nil {
			st = append(st, v)
			continue
		}

		if len(st) == 0 {
			return "", fmt.Errorf("cannot repeat rune: %s", string(v))
		}

		repeatRune := st[len(st)-1]
		st = st[:len(st)-1]

		for i := 0; i < n; i++ {
			st = append(st, repeatRune)
		}
	}

	return string(st), nil
}

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
  - "a4bc2d5e" => "aaaabccddddde"
  - "abcd" => "abcd"
  - "45" => "" (некорректная строка)
  - "" => ""

Дополнительное задание: поддержка escape - последовательностей
  - qwe\4\5 => qwe45 (*)
  - qwe\45 => qwe44444 (*)
  - qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
func main() {
	input := `45`

	result, err := Unpack(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
