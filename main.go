package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	var host string
	var port int

	app := &cli.App{
		Name: "http-echo",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       "127.0.0.1",
				Destination: &host,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       18888,
				Destination: &port,
			},
		},
		Action: func(c *cli.Context) error {
			http.HandleFunc("/", handleHttpEcho)
			addr := fmt.Sprintf("%s:%d", host, port)
			log.Printf("listen %s\n", addr)
			return http.ListenAndServe(addr, nil)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type HttpEchoResponse struct {
	Host    string
	Path    string
	Method  string
	Proto   string
	Query   url.Values
	Headers http.Header
	Body    string
}

func handleHttpEcho(w http.ResponseWriter, r *http.Request) {
	result, err := httpEcho(r)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func httpEcho(r *http.Request) (*HttpEchoResponse, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return &HttpEchoResponse{
		Host:    r.Host,
		Path:    r.URL.Path,
		Method:  r.Method,
		Proto:   r.Proto,
		Query:   r.URL.Query(),
		Headers: r.Header,
		Body:    string(body),
	}, nil
}
