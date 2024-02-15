package pattern

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

// Нужен для добавления поведения для разных объектов без их изменения

type (
	A struct{}
	B struct{}
	C struct{}
)

func (a A) Accept(v Visitor) {
	v.VisitA(a)
}
func (b B) Accept(v Visitor) {
	v.VisitB(b)
}
func (c C) Accept(v Visitor) {
	v.VisitC(c)
}

type Visitable interface {
	Accept(v Visitor)
}

type Visitor interface {
	VisitA(a A)
	VisitB(b B)
	VisitC(c C)
}

// Без добавления кода в структуры
func Example1() {
	v := Visitor(nil)
	sl := []any{}

	for _, e := range sl {
		switch c := e.(type) {
		case A:
			v.VisitA(c)
		case B:
			v.VisitB(c)
		case C:
			v.VisitC(c)
		}
	}

}

// С добавлением
func Example2() {
	v := Visitor(nil)
	sl := []Visitable{}

	for _, e := range sl {
		e.Accept(v)
	}
}
