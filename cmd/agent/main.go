package main

import (
	"log"

	"github.com/iselldonuts/metrics/internal/core"
	"github.com/iselldonuts/metrics/internal/metrics"
)

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	run()
}

func run() {
	log.Printf(
		"Running agent | url: %s, ReportInterval: %d, PollInterval: %d\n",
		options.baseURL, options.reportInterval, options.pollInterval,
	)

	poller := metrics.NewPoller()
	agent := core.NewAgent(poller, core.Config{
		BaseURL:        options.baseURL,
		PollInterval:   options.pollInterval,
		ReportInterval: options.reportInterval,
	})
	agent.Start()
}
