package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/iselldonuts/metrics/internal/config/agent"
)

const defaultReportInterval = 10
const defaultPollInterval = 2

func getConfig() (*agent.Config, error) {
	conf := &agent.Config{}

	flag.StringVar(&conf.Address, "a", "localhost:8080", "Server URL")
	flag.IntVar(&conf.ReportInterval, "r", defaultReportInterval, "Report interval in seconds")
	flag.IntVar(&conf.PollInterval, "p", defaultPollInterval, "Poll interval in seconds")
	flag.Parse()

	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("could not parse env variables: %w", err)
	}

	return conf, nil
}
