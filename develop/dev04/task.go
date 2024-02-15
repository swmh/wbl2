package main

import (
	"fmt"
	"slices"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func GetAnagramLetters(s string) map[rune]int {
	letters := make(map[rune]int)
	for _, v := range s {
		letters[v]++
	}

	return letters
}

func Equal(a map[rune]int, b map[rune]int) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if b[k] != v {
			return false
		}
	}

	return true
}

func GroupAnagrams(words []string) map[string][]string {
	result := make(map[string][]string)
	anagrams := make(map[string]map[rune]int)
	duplicates := make(map[string]bool)

loop:
	for _, word := range words {
		if word == "" || duplicates[word] {
			continue
		}

		duplicates[word] = true

		word = strings.ToLower(word)
		anagram := GetAnagramLetters(word)

		for k := range result {
			if Equal(anagram, anagrams[k]) {
				result[k] = append(result[k], word)
				continue loop
			}
		}

		result[word] = []string{}
		anagrams[word] = anagram
	}

	for k, v := range result {
		if len(v) == 0 {
			delete(result, k)
		}

		if len(v) > 1 {
			slices.Sort(v)
		}

	}

	return result
}

func main() {
	arr := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "den", "слиток"}
	res := GroupAnagrams(arr)
	fmt.Printf("%v\n", res)
}
