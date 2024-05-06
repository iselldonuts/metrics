package main

import "flag"

var options struct {
	baseURL        string
	reportInterval int
	pollInterval   int
}

func parseFlags() {
	flag.StringVar(&options.baseURL, "a", "localhost:8080", "Server URL")
	flag.IntVar(&options.reportInterval, "r", 10, "Report interval in seconds")
	flag.IntVar(&options.pollInterval, "p", 2, "Poll interval in seconds")
	flag.Parse()
}
