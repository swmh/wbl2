package pattern

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

// Зависимость от одного класса вместо нескольких
type Client struct{}

func (c Client) SomeWork() {
	facade := Facade{}
	facade.Operation1()
	facade.Operation2()
	facade.Operation3()
}

type Facade struct {
	Subsystem1
	Subsystem2
	Subsystem3
}

type Subsystem1 struct{}

func (s Subsystem1) Operation1() {}

type Subsystem2 struct{}

func (s Subsystem2) Operation2() {}

type Subsystem3 struct{}

func (s Subsystem2) Operation3() {}
