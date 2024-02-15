package main

import (
	"bufio"
	"io"
	"strconv"
)

type Printer interface {
	Print(Entry)
	Result()
}

type DefaultPrinter struct {
	before, after int
	curAfter      int
	lineNum       bool
	invert        bool
	buf           chan Entry
	b             *bufio.Writer
}

func NewDefaultPrinter(before, after int, w io.Writer, lineNum bool, invert bool) *DefaultPrinter {
	return &DefaultPrinter{
		before:  before,
		after:   after,
		buf:     make(chan Entry, before),
		b:       bufio.NewWriter(w),
		lineNum: lineNum,
		invert:  invert,
	}
}

type Color []byte

var RED = []byte("\033[31m")
var TERMINATE = []byte("\033[0m")

func (d *DefaultPrinter) println(b []byte) {
	d.b.Write(b)
	d.b.WriteByte('\n')
}

func (d *DefaultPrinter) printLineNum(n int, isMatch bool) {
	sep := byte('-')
	if isMatch {
		sep = ':'
	}

	d.b.WriteString(strconv.Itoa(n))
	d.b.WriteByte(sep)
}

func (d *DefaultPrinter) printColorized(b []byte, color Color) {
	d.b.Write(color)
	d.b.Write(b)
	d.b.Write(TERMINATE)
}

func (d *DefaultPrinter) printMatches(b []byte, matches []Match) {
	buf := d.b

	buf.Write(b[:matches[0].start])

	for i, m := range matches {
		if m.len != 0 {
			d.printColorized(b[m.start:m.start+m.len], RED)
		}
		if i+1 < len(matches) {
			buf.Write(b[m.start+m.len : matches[i+1].start])
		}
	}

	buf.Write(b[matches[len(matches)-1].start+matches[len(matches)-1].len:])
	buf.WriteByte('\n')

}

func (d *DefaultPrinter) Print(e Entry) {
	if len(e.matches) == 0 && d.invert {
		for len(d.buf) > 0 {
			ee := <-d.buf

			if d.lineNum {
				d.printLineNum(ee.line, false)
			}

			if len(ee.matches) == 0 {
				d.println(ee.b)
			} else {
				d.printMatches(ee.b, ee.matches)
			}

		}

		if d.lineNum {
			d.printLineNum(e.line, true)
		}

		d.println(e.b)
		d.curAfter = d.after
	}

	if len(e.matches) != 0 {
		for len(d.buf) > 0 {
			ee := <-d.buf

			if d.lineNum {
				d.printLineNum(ee.line, false)
			}
			d.println(ee.b)
		}

		if d.lineNum {
			d.printLineNum(e.line, true)
		}
		d.printMatches(e.b, e.matches)

		d.curAfter = d.after

		return
	}

	if d.curAfter > 0 {
		d.curAfter--

		if d.lineNum {
			d.printLineNum(e.line, false)
		}

		d.println(e.b)
		return
	}

	if d.before != 0 && len(d.buf) >= d.before {
		<-d.buf
	}

	if d.before != 0 {
		d.buf <- e
	}
}

func (d *DefaultPrinter) Result() {
	d.b.Flush()
}
