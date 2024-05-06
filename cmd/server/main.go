package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/api"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", api.UpdateMetric)
	r.Get("/value/{type}/{name}", api.GetMetric)
	r.Get("/", api.Info)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
