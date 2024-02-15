//go:build unix
// +build unix

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

type Ps struct {
	args []string
}

func NewPs(args []string) *Ps {
	return &Ps{args}
}

func (f *Ps) Exec(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintf(stderr, "ps: %s\n", err)
		return
	}

	buf := bufio.NewWriter(stdout)
	defer buf.Flush()

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return
		default:
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()

			if name < "0" || name > "9" {
				continue
			}

			buf.WriteString(name)
			buf.WriteByte('\n')
		}
	}
}
