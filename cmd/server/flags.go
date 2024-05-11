package main

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

var baseURL string

type envs struct {
	Address string `env:"ADDRESS"`
}

func parseFlags() error {
	flag.StringVar(&baseURL, "a", "localhost:8080", "Server URL")
	flag.Parse()

	cfg := envs{}
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("could not parse env variables: %w", err)
	}

	if cfg.Address != "" {
		baseURL = cfg.Address
	}
	return nil
}
