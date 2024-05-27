package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/api"
	"github.com/iselldonuts/metrics/internal/config/server"
	"github.com/iselldonuts/metrics/internal/middleware"
	"github.com/iselldonuts/metrics/internal/storage"
	"github.com/iselldonuts/metrics/internal/storage/memory"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic("cannot sync zap")
		}
	}(logger)
	log := logger.Sugar()

	conf, err := getConfig()
	if err != nil {
		log.Panic(err)
	}

	if err := run(conf, log); err != nil {
		log.Panic(err)
	}
}

func run(conf *server.Config, log *zap.SugaredLogger) error {
	r := chi.NewRouter()
	s := storage.NewStorage(storage.Config{
		Memory: &memory.Config{},
	})

	r.Use(middleware.Logger(log))
	r.Use(middleware.Gzip(log))

	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric(s))
	r.Post("/update/", api.UpdateMetricJSON(s))
	r.Get("/value/{type}/{name}", api.GetMetric(s))
	r.Post("/value/", api.GetMetricJSON(s))
	r.Get("/", api.Info(s))

	log.Infow("Running server", "url", conf.Address)

	if err := http.ListenAndServe(conf.Address, r); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server run error: %w", err)
		}
	}
	return nil
}
