package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

type Productt interface {
	Do()
}

type ProductA struct {
	f bool
}

func (p *ProductA) Do() {
	fmt.Println(p.f)
}

type ProductB struct {
	f string
}

func (p *ProductB) Do() {
	fmt.Println(p.f)
}

type Creator interface {
	CreateProduct(name string) Productt
}

type ConcreteCreator struct {
}

func (c *ConcreteCreator) CreateProduct(name string) Productt {
	switch name {
	case "A":
		return &ProductA{f: true}
	case "B":
		return &ProductB{f: "Hello"}
	default:
		return nil
	}
}
