package main

import (
	"github.com/swmh/wbl2/develop/dev01/lib"
)

//func PrintTime() {
//	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
//	err = fmt.Errorf("aboba")
//	if err != nil {
//		log.Fatalln(err)
//	}
//
//	ntpTime := time.Now().Add(response.ClockOffset)
//	now := time.Now()
//
//	fmt.Printf("%v / %v\n", now, ntpTime)
//}
//
/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func main() {
	lib.PrintTime()
}
