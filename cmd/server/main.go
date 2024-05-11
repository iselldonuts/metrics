package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/api"
)

func main() {
	if err := parseFlags(); err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric)
	r.Get("/value/{type}/{name}", api.GetMetric)
	r.Get("/", api.Info)

	log.Println("Running server on", baseURL)
	return fmt.Errorf("run error: %w", http.ListenAndServe(baseURL, r))
}
