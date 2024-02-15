package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Config struct {
	After  int
	Before int

	NumberOfMatches bool

	IgnoreCase bool
	Invert  bool
	Fixed   bool
	LineNum bool

	Inputs  []string
	Pattern string
}

func NewConfig() (*Config, error) {
	c := &Config{}

	var before int
	var after int
	var cont int

	flag.IntVar(&cont, "C", 0, "")
	flag.IntVar(&before, "B", 0, "")
	flag.IntVar(&after, "A", 0, "")

	flag.BoolVar(&c.NumberOfMatches, "c", false, "")
	flag.BoolVar(&c.IgnoreCase, "i", false, "")
	flag.BoolVar(&c.Invert, "v", false, "")
	flag.BoolVar(&c.Fixed, "F", false, "")
	flag.BoolVar(&c.LineNum, "n", false, "")

	flag.Parse()

	if cont < 0 {
		fmt.Println("invalid context length argument")
		os.Exit(1)
	}

	c.Before = cont
	c.After = cont

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "A":
			c.After = after
		case "B":
			c.Before = before
		}
	})

	pattern := flag.Arg(0)

	for _, v := range flag.Args()[1:] {
		c.Inputs = append(c.Inputs, v)
	}

	c.Pattern = pattern

	return c, nil
}

func (c *Config) GetPrinter(w io.Writer) (Printer, error) {
	if c.NumberOfMatches {
		return NewPrinterCounter(w), nil
	}

	return NewDefaultPrinter(c.Before, c.After, w, c.LineNum, c.Invert), nil
}

func (c *Config) GetMatcher() (Matcher, error) {
	return NewRegexpMatcher(c.Pattern, c.Fixed, c.IgnoreCase)

}
