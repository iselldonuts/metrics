package main

import (
	"log"

	"github.com/iselldonuts/metrics/internal/config/agent"
	"github.com/iselldonuts/metrics/internal/core"
	"github.com/iselldonuts/metrics/internal/metrics"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	run(cfg)
}

func run(conf *agent.Config) {
	log.Printf(
		"Running agent | url: %s, ReportInterval: %d, PollInterval: %d\n",
		conf.Address, conf.ReportInterval, conf.PollInterval,
	)

	poller := metrics.NewPoller()
	a := core.NewAgent(poller, conf)
	a.Start()
}
