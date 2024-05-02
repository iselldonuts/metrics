package api

import (
	"github.com/iselldonuts/metrics/internal/storage"
	"github.com/iselldonuts/metrics/internal/storage/memory"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Storage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value float64)
}

var s = storage.NewStorage(storage.Config{
	Memory: &memory.Config{},
})

func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/update/")
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		if len(parts) == 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		return
	}

	metricType, name, value := parts[0], parts[1], parts[2]
	switch metricType {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, "wrong metric value", http.StatusBadRequest)
		}

		s.UpdateGauge(name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, "wrong metric value", http.StatusBadRequest)
		}

		s.UpdateCounter(name, v)
	default:
		http.Error(w, "wrong metric type", http.StatusBadRequest)
		return
	}
}

func Info(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, s.String())
}
