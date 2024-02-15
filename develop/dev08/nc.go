package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type NC struct {
	err     error
	listen  bool
	address string
	network string
	bufSize int
}

func NewNC(args []string) *NC {
	nc := NC{
		network: TCP,
		bufSize: 512,
	}

	var isUDP bool

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.BoolVar(&nc.listen, "l", false, "listen")
	fs.BoolVar(&isUDP, "u", false, "listen")

	err := fs.Parse(args)
	if err != nil {
		return &NC{
			err: err,
		}
	}

	if isUDP {
		nc.network = UDP
	}

	if nc.listen {
		port := fs.Arg(0)
		if port == "" {
			return &NC{
				err: errors.New("no port specified"),
			}
		}
		nc.address = fmt.Sprintf(":%s", fs.Arg(0))
	} else {
		address := fs.Arg(0)
		if address == "" {
			return &NC{
				err: errors.New("no address specified"),
			}
		}

		port := fs.Arg(1)
		if port == "" {
			return &NC{
				err: errors.New("no port specified"),
			}
		}

		nc.address = fmt.Sprintf("%s:%s", address, port)
	}

	fmt.Printf("%+v\n", nc)

	return &nc
}

func (nc *NC) ListenConn(ctx context.Context, listener io.Reader, stdout, stderr io.Writer) {
	buf := make([]byte, nc.bufSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := listener.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				fmt.Fprintf(stderr, "nc: %s\n", err)
				return
			}

			stdout.Write(buf[:n])
		}
	}
}

func (nc *NC) Listen(ctx context.Context, stdout, stderr io.Writer) error {
	switch nc.network {
	case TCP:
		ln, err := net.Listen(TCP, nc.address)
		if err != nil {
			return err
		}

		go func() {
			<-ctx.Done()
			ln.Close()
		}()

		for {
			conn, err := ln.Accept()
			if err != nil {
				return err
			}

			go func() {
				nc.ListenConn(ctx, conn, stdout, stderr)
				conn.Close()
			}()
		}

	case UDP:
		laddr, err := net.ResolveUDPAddr(UDP, nc.address)
		if err != nil {
			return err
		}

		listener, err := net.ListenUDP(UDP, laddr)
		if err != nil {
			return err
		}

		nc.ListenConn(ctx, listener, stdout, stderr)
		listener.Close()
	default:
		return fmt.Errorf("unknown network %s", nc.network)
	}

	return nil
}

func (nc *NC) Connect(ctx context.Context, stdin io.Reader, stderr io.Writer) {
	conn, err := net.Dial(nc.network, nc.address)
	if err != nil {
		fmt.Fprintf(stderr, "nc: %s\n", err)
		return
	}

	defer conn.Close()

	ch := make(chan string, 1)
	defer close(ch)

	go func() {
		buf := bufio.NewReader(stdin)
		for {
			v, err := buf.ReadString('\n')
			if err != nil {
				return
			}

			ch <- v
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-ch:
			_, err = conn.Write([]byte(v))
			if err != nil {
				fmt.Fprintf(stderr, "nc: %s\n", err)
				return
			}
		}
	}
}

func (nc *NC) Exec(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) {
	if nc.err != nil {
		fmt.Fprintf(stderr, "nc: %s\n", nc.err)
		return
	}

	if nc.listen {
		err := nc.Listen(ctx, stdout, stderr)
		if err != nil {
			fmt.Fprintf(stderr, "nc: %s\n", err)
			return
		}

		return

	}

	nc.Connect(ctx, stdin, stderr)
}
