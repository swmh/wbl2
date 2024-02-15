package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

// Работает как паттерн стратегия, но в отличии от него может менять состояние контекста

type Contextt struct {
	state State
}

func (c *Contextt) ChangeState(s State) {
	c.state = s
}
func (c *Contextt) Do() {
	c.state.Do()
}

type State interface {
	Do()
}

type StateA struct {
	context *Contextt
}

func (s *StateA) Do() {
	fmt.Println("StateA")
	s.context.ChangeState(&StateB{s.context})
}

type StateB struct {
	context *Contextt
}

func (s *StateB) Do() {
	fmt.Println("StateB")
}

func Example3() {
	context := &Contextt{}
	context.ChangeState(&StateA{context})
	context.Do()
	context.ChangeState(&StateB{context})
	context.Do()
}
