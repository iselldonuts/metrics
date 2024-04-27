package main

import (
	"github.com/iselldonuts/metrics/internal"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	storage := internal.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		}

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/update/"), "/")
		if len(parts) < 3 {
			if len(parts) == 1 {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
		}

		metricType, name, value := parts[0], parts[1], parts[2]
		switch metricType {
		case "gauge":
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
			storage.Gauge[name] = v
		case "counter":
			v, err := strconv.Atoi(value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
			storage.Counter[name] += int64(v)
		default:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	})

	mux.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, storage.String())
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
