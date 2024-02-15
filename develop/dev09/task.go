package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Config struct {
	Client     http.Client
	Url        *url.URL
	Link       string
	Downloaded *sync.Map
	I          *atomic.Int64

	Depth     int
	Recursive bool
}

func Parse() (Config, error) {
	cfg := Config{
		Recursive:  false,
		Url:        &url.URL{},
		Link:       "",
		Client:     http.Client{},
		Downloaded: &sync.Map{},
	}

	flag.IntVar(&cfg.Depth, "d", 5, "depth")
	flag.BoolVar(&cfg.Recursive, "r", false, "recursive")

	flag.Parse()

	link := flag.Arg(0)
	if link == "" {
		return cfg, errors.New("no link specified")
	}

	l, err := url.Parse(link)
	if err != nil {
		return cfg, err
	}

	cfg.Url = l
	cfg.Link = link
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
	}
	cfg.Client = http.Client{
		Transport: transport,
		Jar:       nil,
		Timeout:   0,
	}

	return cfg, nil
}

// regexp: href|src\s*=\s*["\']([^"\'\s>]+)["\']?

// better <a.+?\s*href\s*=\s*["\']?([^\#"\'\s>]+)["\']?
var (
	href = regexp.MustCompile(`<a.+?\s*href\s*=\s*["\']?([^\#"\'\s>]+)["\']?`)
)

func (c *Config) GetLinks(b []byte) []string {
	var result []string

	matches := href.FindAllSubmatch(b, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		for _, v := range match[1:] {
			u, err := url.Parse(string(v))
			if err != nil {
				continue
			}

			if u.IsAbs() {
				if u.Host != c.Url.Host {
					continue
				}
			}

			u = c.Url.ResolveReference(u)
			result = append(result, u.String())
		}
	}

	return result
}

func ResolveURL(u string) string {
	if strings.HasSuffix(u, "/") {
		u = u + "index.html"
	}

	return u
}

func (cfg *Config) Get(link string) (io.ReadCloser, error) {
	resp, err := cfg.Client.Get(link)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func GetName(s string) string {
	var name string
	parts := strings.Split(s, "/")

	name = parts[len(parts)-1]
	if name == "" {
		name = "index.html"
	}

	return name
}

func GetPathFromURL(u string) (string, error) {
	uu, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	r := path.Join(uu.Host, uu.Path)

	return r, nil
}

func MkDir(name string) error {
	err := os.MkdirAll(name, 0777)
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Err.Error() == "file exists" {
			return nil
		}

		return err
	}

	return nil
}

func (c *Config) RecuresiveGet(link string, level int) error {
	_, ok := c.Downloaded.LoadOrStore(link, struct{}{})
	if ok {
		c.I.Add(1)
		return nil
	}

	fmt.Printf("Downloading %s\n", link)
	resp, err := c.Get(link)
	if err != nil {
		return err
	}
	defer resp.Close()

	link = ResolveURL(link)

	filePath, err := GetPathFromURL(link)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath = path.Join(cwd, filePath)
	fmt.Println("Save to ", filePath)

	dir, _ := path.Split(filePath)
	if err = MkDir(dir); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var b bytes.Buffer
	wr := io.MultiWriter(&b, file)

	_, err = io.Copy(wr, resp)
	if err != nil {
		return err
	}

	links := c.GetLinks(b.Bytes())

	if level < c.Depth {
		wg := sync.WaitGroup{}
		wg.Add(len(links))

		for _, v := range links {
			go func(u string) {
				defer wg.Done()
				if err = c.RecuresiveGet(u, level+1); err != nil {
					fmt.Printf("Fail %s\n", err)
				}
			}(v)
		}

		wg.Wait()
	}

	return nil
}

func main() {
	cfg, err := Parse()
	if err != nil {
		log.Fatalln(err)
	}

	if !cfg.Recursive {
		resp, err := cfg.Get(cfg.Url.String())
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Close()

		name := GetName(cfg.Link)

		file, err := os.Create(name)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		io.Copy(file, resp)
		return
	}

	err = cfg.RecuresiveGet(cfg.Url.String(), 0)
	if err != nil {
		fmt.Printf("Fail %s\n", err)
	}
	fmt.Println(cfg.I.Load())
}
