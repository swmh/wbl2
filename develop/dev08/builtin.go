package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// встроенные команды: cd/pwd/echo/kill/ps
// поддержать fork/exec команды
// конвеер на пайпах
//
// Реализовать утилиту netcat (nc) клиент
// принимать данные из stdin и отправлять в соединение (tcp/udp)

func GetExecutablePath(s string) string {
	if strings.ContainsRune(s, os.PathSeparator) {
		return s
	}

	r, _ := exec.LookPath(s)
	return r
}

type CdCommand struct {
	dir string
}

func NewCd(args []string) *CdCommand {
	return &CdCommand{args[0]}
}

func (c *CdCommand) Exec(_ context.Context, _ io.Reader, _, stderr io.Writer) {
	err := os.Chdir(string(c.dir))
	if err != nil {
		fmt.Fprintf(stderr, "cd: %s", err)
	}
}

type PwdCommand struct{}

func NewPwd(_ []string) *PwdCommand {
	return &PwdCommand{}
}

func (p *PwdCommand) Exec(_ context.Context, _ io.Reader, stdout, stderr io.Writer) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "pwd: %s", err)
		return
	}

	fmt.Fprintln(stdout, dir)
}

type EchoCommand struct {
	text string
}

func NewEcho(args []string) *EchoCommand {
	return &EchoCommand{strings.Join(args, " ")}
}

func (e *EchoCommand) Exec(_ context.Context, _ io.Reader, stdout, _ io.Writer) {
	fmt.Fprintln(stdout, e.text)
}

type Exec struct {
	args []string
}

func NewExec(args []string) *Exec {
	return &Exec{args}
}

func (f *Exec) Exec(_ context.Context, _ io.Reader, _, stderr io.Writer) {
	argv0 := GetExecutablePath(f.args[0])
	err := syscall.Exec(argv0, f.args, os.Environ())
	if err != nil {
		fmt.Fprintf(stderr, "exec: %s\n", err)
	}
}

type Fork struct {
	args []string
}

func NewFork(args []string) *Fork {
	return &Fork{args}
}

func (f *Fork) Exec(_ context.Context, _ io.Reader, _, stderr io.Writer) {
	argv0 := GetExecutablePath(f.args[0])
	attr := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{},
		Sys: &syscall.SysProcAttr{
			Foreground: false,
		},
	}
	_, err := syscall.ForkExec(argv0, f.args, &attr)
	if err != nil {
		fmt.Fprintf(stderr, "fork: %s", err)
	}
}
