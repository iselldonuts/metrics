package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/storage"
	"github.com/iselldonuts/metrics/internal/storage/memory"
)

var s = storage.NewStorage(storage.Config{
	Memory: &memory.Config{},
})

func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	switch mtype {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "wrong metric value", http.StatusBadRequest)
			return
		}

		s.UpdateGauge(name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "wrong metric value", http.StatusBadRequest)
			return
		}

		s.UpdateCounter(name, v)
	default:
		http.Error(w, "wrong metric type", http.StatusBadRequest)
		return
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	switch mtype {
	case "gauge":
		m, ok := s.GetGauge(name)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprint(w, m)
	case "counter":
		m, ok := s.GetCounter(name)
		if !ok {
			http.Error(w, "metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprint(w, m)
	default:
		http.Error(w, "wrong metric type", http.StatusBadRequest)
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	_, _ = fmt.Fprintln(w, "<p>Counter metrics:</p><ul>")
	for name, value := range s.GetAllCounter() {
		_, _ = fmt.Fprintf(w, "<li>%s: %v</li>", name, value)
	}
	_, _ = fmt.Fprintln(w, "</ul><p>Gauge metrics:</p><ul>")
	for name, value := range s.GetAllGauge() {
		_, _ = fmt.Fprintf(w, "<li>%s: %v</li>", name, value)
	}
	_, _ = fmt.Fprintln(w, "</ul>")
}
