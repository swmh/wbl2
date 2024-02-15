package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.

*/

type Config struct {
	host    string
	port    int
	timeout time.Duration
}

func Parse() Config {
	var c Config

	flag.DurationVar(&c.timeout, "timeout", 10*time.Second, "timeout")

	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		log.Fatalf("hostname not specified\n")
	}

	port := flag.Arg(1)
	if port == "" {
		log.Fatalf("port not specified\n")
	}

	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("port must be a number")
	}

	c.host = host
	c.port = p

	return c
}

func Connect(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}

	buf := make([]byte, 512)

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return err
		}

		conn.SetWriteDeadline(time.Now().Add(time.Second))
		_, err = conn.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	cfg := Parse()
	err := Connect(cfg.host+":"+strconv.Itoa(cfg.port), cfg.timeout)
	if !errors.Is(err, io.EOF) {
		fmt.Println(err)
	}
}
