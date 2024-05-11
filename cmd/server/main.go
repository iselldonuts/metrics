package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/api"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric)
	r.Get("/value/{type}/{name}", api.GetMetric)
	r.Get("/", api.Info)

	fmt.Println("Running server on", baseURL)
	return fmt.Errorf("run error: %w", http.ListenAndServe(baseURL, r))
}
