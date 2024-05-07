package main

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"log"
)

var baseURL string

type envs struct {
	Address string `env:"ADDRESS"`
}

func parseFlags() {
	flag.StringVar(&baseURL, "a", "localhost:8080", "Server URL")
	flag.Parse()

	cfg := envs{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Address != "" {
		baseURL = cfg.Address
	}
}
