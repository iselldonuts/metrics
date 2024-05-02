package main

import (
	"github.com/iselldonuts/metrics/internal/api"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", api.UpdateMetric)
	mux.HandleFunc("/mem", api.Info)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
