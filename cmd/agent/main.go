package main

import (
	"log"

	"github.com/iselldonuts/metrics/internal/core"
	"github.com/iselldonuts/metrics/internal/metrics"
)

func main() {
	parseFlags()
	run()
}

func run() {
	log.Printf(
		"Running agent | url: %s, ReportInterval: %d, PollInterval: %d\n",
		options.baseURL, options.reportInterval, options.pollInterval,
	)

	poller := metrics.NewPoller()

	a := core.NewAgent(poller, core.Config{
		BaseURL:        options.baseURL,
		PollInterval:   options.pollInterval,
		ReportInterval: options.reportInterval,
	})
	a.Start()
}
