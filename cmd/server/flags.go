package main

import (
	"flag"
)

var baseURL string

func parseFlags() {
	flag.StringVar(&baseURL, "a", "localhost:8080", "Server URL")
	flag.Parse()
}
