package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/model"
	"github.com/iselldonuts/metrics/internal/storage"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
	Gauge           = "gauge"
	Counter         = "counter"
)

var (
	errMetricType     = errors.New("wrong metric type")
	errMetricValue    = errors.New("wrong metric value")
	errMetricNotFound = errors.New("metric not found")
)

func UpdateMetric(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mtype := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")
		value := chi.URLParam(r, "value")

		switch mtype {
		case Gauge:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(w, errMetricValue.Error(), http.StatusBadRequest)
				return
			}
			s.UpdateGauge(name, v)
		case Counter:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				http.Error(w, errMetricValue.Error(), http.StatusBadRequest)
				return
			}
			s.UpdateCounter(name, v)
		default:
			http.Error(w, errMetricType.Error(), http.StatusBadRequest)
			return
		}
	}
}

func UpdateMetricJSON(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(ContentType) != ContentTypeJSON {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		var metrics model.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			log.Printf("failed to unmarshal metrics: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch metrics.MType {
		case Gauge:
			s.UpdateGauge(metrics.ID, *metrics.Value)
		case Counter:
			s.UpdateCounter(metrics.ID, *metrics.Delta)
		default:
			http.Error(w, errMetricType.Error(), http.StatusBadRequest)
			return
		}
	}
}

func GetMetric(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mtype := chi.URLParam(r, "type")
		name := chi.URLParam(r, "name")

		switch mtype {
		case Gauge:
			v, ok := s.GetGauge(name)
			if !ok {
				http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set(ContentType, "text/plain")

			value := strconv.FormatFloat(v, 'f', -1, 64)
			_, err := w.Write([]byte(value))
			if err != nil {
				log.Printf("error writing response: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		case Counter:
			v, ok := s.GetCounter(name)
			if !ok {
				http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set(ContentType, "text/plain")

			value := strconv.FormatInt(v, 10)
			_, err := w.Write([]byte(value))
			if err != nil {
				log.Printf("error writing response: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, errMetricType.Error(), http.StatusBadRequest)
		}
	}
}

func GetMetricJSON(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(ContentType) != ContentTypeJSON {
			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
			return
		}

		var metrics model.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set(ContentType, ContentTypeJSON)

		switch metrics.MType {
		case Gauge:
			v, ok := s.GetGauge(metrics.ID)
			if !ok {
				http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
				return
			}

			metrics.Value = &v
			if err := json.NewEncoder(w).Encode(&metrics); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case Counter:
			v, ok := s.GetCounter(metrics.ID)
			if !ok {
				http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
				return
			}
			metrics.Delta = &v

			encoder := json.NewEncoder(w)
			if err := encoder.Encode(&metrics); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			http.Error(w, errMetricType.Error(), http.StatusBadRequest)
		}
	}
}

func Info(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(ContentType, "text/html")

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
}
