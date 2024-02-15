package pattern

/*
	Реализовать паттерн «цепочка вызовов».

Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.

	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/
type request int

type Handler interface {
	SetNext(h Handler)
	Handle(r request)
}

type ConcreteHandler struct {
	next Handler
}

func (c *ConcreteHandler) SetNext(h Handler) {
	c.next = h
}

func (c *ConcreteHandler) Handle(r request) {
	if r == 0 {
		if c.next != nil {
			c.next.Handle(r)
		}
		return
	}

	println("Valid")
}
