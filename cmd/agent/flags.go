package main

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

var options struct {
	baseURL        string
	reportInterval int
	pollInterval   int
}

type envs struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func parseFlags() {
	flag.StringVar(&options.baseURL, "a", "localhost:8080", "Server URL")
	flag.IntVar(&options.reportInterval, "r", 10, "Report interval in seconds")
	flag.IntVar(&options.pollInterval, "p", 2, "Poll interval in seconds")
	flag.Parse()

	cfg := envs{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address != "" {
		options.baseURL = cfg.Address
	}
	if cfg.ReportInterval != 0 {
		options.reportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		options.pollInterval = cfg.PollInterval
	}
}
