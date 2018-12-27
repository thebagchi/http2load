package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type HTTPRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Queries map[string][]string `json:"queries"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
	Expect  int                 `json:"expect"`
}

var client = http.Client{
	Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (conn net.Conn, e error) {
			return net.Dial(network, addr)
		},
	},
}

/*
func init() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	h2s := &http2.Server{
		// ...
	}
	h1s := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handler, h2s),
	}
	go h1s.ListenAndServe()
}
*/

func main() {
	var (
		script = flag.String("script", "request.star", "configuration script filename")
	)
	flag.Parse()

	if nil == script {
		log.Fatal("No configuration script provided, exiting...")
	}

	if len(strings.TrimSpace(*script)) == 0 {
		log.Fatal("Configuration script filename cannot be empty")
	}

	repeat, requests := LoadScript(*script)

	if repeat > 0 {
		gauge := NewGauge()
		var values []HTTPRequest
		err := json.Unmarshal([]byte(requests), &values)
		if nil == err {
			pool := NewWorkerPool(8)
			for count := 0; count < repeat; count++ {
				for _, request := range values {
					request := request
					pool.Enqueue(func() {
						start := time.Now()
						req, err := http.NewRequest(request.Method, request.Path, strings.NewReader(request.Body))
						if nil == err {
							res, err := client.Do(req)
							if nil != res {
								defer res.Body.Close()
							}
							if nil == err {
								if request.Expect != 0 {
									if request.Expect != res.StatusCode {
										fmt.Println(fmt.Sprintf(
											"Request didnot complete with expected status code %d received %d",
											request.Expect,
											res.StatusCode,
										))
									}
								}
								fmt.Println(res.StatusCode)
							} else {
								fmt.Println("Error: ", err)
							}
						} else {
							fmt.Println("Error: ", err)
						}
						diff := time.Now().Sub(start).Nanoseconds() / (1000 * 1000)
						gauge.Add(float64(diff))
					})
				}
			}
			pool.Await()
			fmt.Println(fmt.Sprintf(
				"Min: %v(ms) Max: %v(ms) Mean: %v(ms) Total: %v",
				gauge.Min,
				gauge.Max,
				gauge.Mean(),
				gauge.Count,
			))
		}
	}
}
