package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/api"
	"github.com/iselldonuts/metrics/internal/config/server"
	"github.com/iselldonuts/metrics/internal/storage"
	"github.com/iselldonuts/metrics/internal/storage/memory"
)

func main() {
	conf, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(conf); err != nil {
		log.Fatal(err)
	}
}

func run(conf *server.Config) error {
	r := chi.NewRouter()
	s := storage.NewStorage(storage.Config{
		Memory: &memory.Config{},
	})

	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric(s))
	r.Get("/value/{type}/{name}", api.GetMetric(s))
	r.Get("/", api.Info(s))

	log.Printf("Running server on %s", conf.Address)

	if err := http.ListenAndServe(conf.Address, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server run error: %w", err)
		}
	}
	return nil
}
