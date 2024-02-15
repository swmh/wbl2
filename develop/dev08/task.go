package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Command interface {
	Exec(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer)
}

var builtinCommands = map[string]func(args []string) Command{
	"cd":   func(args []string) Command { return NewCd(args) },
	"pwd":  func(args []string) Command { return NewPwd(args) },
	"echo": func(args []string) Command { return NewEcho(args) },
	"kill": func(args []string) Command { return NewCd(args) },
	"ps":   func(args []string) Command { return NewPs(args) },
	"fork": func(args []string) Command { return NewFork(args) },
	"exec": func(args []string) Command { return NewExec(args) },
	"nc":   func(args []string) Command { return NewNC(args) },
}

func GetCommand(name string, args []string) Command {
	var command Command
	if builtin, ok := builtinCommands[name]; ok {
		command = builtin(args)
	} else {
		command = &ExternalCommand{
			name: name,
			args: args,
		}
	}

	return command
}

func ParseArgs(s string) []string {
	s = strings.TrimSpace(s)
	args := strings.Split(s, " ")
	for i, v := range args {
		args[i] = strings.Trim(v, `'"`)
	}

	return args
}

type Shell struct {
	prompt string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	jobs   []context.CancelFunc
}

func (s *Shell) Prompt() string {
	return s.prompt
}

func (s *Shell) AddJob(cancel context.CancelFunc) {
	s.jobs = append(s.jobs, cancel)
}

func (s *Shell) StopJobs() {
	for _, cancel := range s.jobs {
		cancel()
	}
}

func main() {
	shell := &Shell{
		prompt: "> ",
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

		for {
			<-sigChan
			fmt.Println("tsd")
			shell.StopJobs()
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprintf(shell.stdout, "\n%s", shell.Prompt())

		if !scanner.Scan() {
			continue
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "exit" {
			os.Exit(0)
		}

		var command Command

		parts := strings.Split(input, "|")

		if len(parts) > 1 {
			st := make([]Command, 0, 2)

			for _, part := range parts {
				args := ParseArgs(part)
				c := GetCommand(args[0], args[1:])

				command = c

				if len(st) == 0 {
					st = append(st, c)
					continue
				}

				p := NewPipe(st[len(st)-1], c)
				st = st[:len(st)-1]

				command = p
			}

		} else {
			args := ParseArgs(input)
			command = GetCommand(args[0], args[1:])
		}

		ctx, cancel := context.WithCancel(context.Background())

		shell.AddJob(cancel)

		command.Exec(ctx, shell.stdin, shell.stdout, shell.stderr)
	}
}
