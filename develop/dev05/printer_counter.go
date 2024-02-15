package main

import (
	"fmt"
	"io"
)

type PrinterCounter struct {
	counter int
	w       io.Writer
}

func NewPrinterCounter(w io.Writer) *PrinterCounter {
	return &PrinterCounter{
		w: w,
	}
}

func (p *PrinterCounter) Print(entry Entry) {
	if len(entry.matches) > 0 {
		p.counter++
	}
}

func (p PrinterCounter) Result() {
	fmt.Fprintln(p.w, p.counter)
}
