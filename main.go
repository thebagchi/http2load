package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
)

type HTTPRequest struct {
	Path    string              `json:"path"`
	Queries map[string][]string `json:"queries"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

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
		var values []HTTPRequest
		err := json.Unmarshal([]byte(requests), &values)
		if nil == err {
			pool := NewWorkerPool(8)
			for count := 0; count < repeat; count++ {
				for _, request := range values {
					request := request
					pool.Enqueue(func() {
						fmt.Println(request)
					})
				}
			}
			pool.Await()
		}
	}
}
