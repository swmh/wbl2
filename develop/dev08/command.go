package main

import (
	"context"
	"io"
	"os/exec"
	"sync"
)

type ExternalCommand struct {
	name string
	args []string
}

func (ec *ExternalCommand) Exec(_ context.Context, stdin io.Reader, stdout, stderr io.Writer) {
	cmd := exec.Command(ec.name, ec.args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	cmd.Run()
}

type Pipe struct {
	c1 Command
	c2 Command
}

func (p *Pipe) Exec(ctx context.Context, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	var wg sync.WaitGroup
	rp, wp := io.Pipe()

	wg.Add(1)
	go func() {
		p.c1.Exec(ctx, stdin, wp, stderr)
		wp.Close()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		p.c2.Exec(ctx, rp, stdout, stderr)
		rp.Close()
		wg.Done()
	}()

	wg.Wait()
}

func NewPipe(c1, c2 Command) *Pipe {
	return &Pipe{
		c1: c1,
		c2: c2,
	}
}
