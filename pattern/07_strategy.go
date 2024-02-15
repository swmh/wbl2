package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

// Определяет схожие алгоритмы и выносит их в отдельные классы

type Strategy interface {
	Do(a, b int)
}

type Context struct {
	strategy Strategy
}

func (c *Context) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

func (c *Context) DoSomething() {
	c.strategy.Do(1, 2)
}
