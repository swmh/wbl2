package pattern

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

// Отделяет действие от ее вызова

type Command interface {
	Execute()
	Undo()
}

type Receiver struct {
}

func (r *Receiver) Action() {
	println("Action")
}

type ConcreteCommand struct {
	receiver Receiver
}

func NewConcreteCommand(receiver Receiver) *ConcreteCommand {
	return &ConcreteCommand{receiver}
}

func (c *ConcreteCommand) Execute() {
	c.receiver.Action()
}

func (c *ConcreteCommand) Undo() {
}

type Invoker struct {
	command Command
}

func (i *Invoker) SetCommand(command Command) {
	i.command = command
}

func (i *Invoker) Invoke() {
	i.command.Execute()
}

func (i *Invoker) Cancel() {
	i.command.Undo()
}

func Example() {
	receiver := Receiver{}
	command := NewConcreteCommand(receiver)
	invoker := Invoker{}
	invoker.SetCommand(command)
	invoker.Invoke()
	invoker.Cancel()
}
