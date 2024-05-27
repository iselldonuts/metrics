package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/iselldonuts/metrics/internal/config/server"
)

func getConfig() (*server.Config, error) {
	conf := &server.Config{}

	flag.StringVar(&conf.Address, "a", "localhost:8080", "Server URL")
	flag.Parse()

	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("could not parse env variables: %w", err)
	}

	return conf, nil
}
