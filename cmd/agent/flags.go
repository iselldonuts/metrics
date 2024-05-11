package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
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

func parseFlags() error {
	const defaultReportInterval = 10
	const defaultPollInterval = 2

	flag.StringVar(&options.baseURL, "a", "localhost:8080", "Server URL")
	flag.IntVar(&options.reportInterval, "r", defaultReportInterval, "Report interval in seconds")
	flag.IntVar(&options.pollInterval, "p", defaultPollInterval, "Poll interval in seconds")
	flag.Parse()

	cfg := envs{}
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("could not parse env variables: %w", err)
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
	return nil
}
