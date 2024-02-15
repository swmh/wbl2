package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

// Нужен для инкапсуляции и разделения создания объекта на отдельные этапы

type Product struct {
	a, b, c int
}

type Builder interface {
	BuildPartA()
	BuildPartB()
	BuildPartC()
	GetResult() Product
}

type Director struct {
	builder Builder
}

func NewDirector(b Builder) Director {
	return Director{builder: b}
}

func (d *Director) Construct() {
	d.builder.BuildPartA()
	d.builder.BuildPartB()
	d.builder.BuildPartC()
}

type ConcreteBuilder struct {
	product Product
}

func NewConcreteBuilder() *ConcreteBuilder {
	return &ConcreteBuilder{}
}

func (b *ConcreteBuilder) BuildPartA() {
	b.product.a = 1
}
func (b *ConcreteBuilder) BuildPartB() {
	b.product.b = 2
}
func (b *ConcreteBuilder) BuildPartC() {
	b.product.c = 3
}
func (b *ConcreteBuilder) GetResult() Product {
	return b.product
}

func builder() {
	b := NewConcreteBuilder()
	d := NewDirector(b)
	d.Construct()
	p := b.GetResult()
	fmt.Printf("a: %d, b: %d, c: %d\n", p.a, p.b, p.c)
}
