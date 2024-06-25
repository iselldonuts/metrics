package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/api"
	"github.com/iselldonuts/metrics/internal/config/server"
	"github.com/iselldonuts/metrics/internal/middleware"
	"github.com/iselldonuts/metrics/internal/storage"
	"github.com/iselldonuts/metrics/internal/storage/file"
	"github.com/iselldonuts/metrics/internal/storage/memory"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
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
	var c storage.Config
	if conf.FileStoragePath == "" {
		c.Memory = &memory.Config{}
	} else {
		c.File = &file.Config{Path: conf.FileStoragePath}
	}

	syncSave := conf.StoreInterval == 0
	s := storage.NewStorage(c, log)
	if conf.FileStoragePath != "" {
		if conf.Restore {
			if err := s.Load(); err != nil {
				log.Infof("Error loading metrics from %q: %v", conf.FileStoragePath, err)
			}
		}

		if !syncSave {
			go func() {
				storeSaveTicker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
				for {
					<-storeSaveTicker.C
					_ = s.Save()
				}
			}()
		}
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(log))
	r.Use(middleware.Gzip(gzip.NewWriter(nil), log))

	a := api.NewAPI(s, log, syncSave)
	r.Mount("/", a.Routes())

	log.Infof("Running server: %+v", *conf)
	if err := http.ListenAndServe(conf.Address, r); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			if err := s.Save(); err != nil {
				return fmt.Errorf("cannot save metrics to disk: %w", err)
			}
		}
		return fmt.Errorf("http server run error: %w", err)
	}
	return nil
}
