package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iselldonuts/metrics/internal/model"
	"github.com/iselldonuts/metrics/internal/storage"
)

const (
	ContentType       = "Content-Type"
	ContentEncoding   = "Content-Encoding"
	AcceptEncoding    = "Accept-Encoding"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html"
	Gauge             = "gauge"
	Counter           = "counter"
	InvalidMetricType = "invalid metric type: %s"
)

var (
	errMetricType     = errors.New("wrong metric type")
	errMetricValue    = errors.New("wrong metric value")
	errMetricNotFound = errors.New("metric not found")
)

type Logger interface {
	Infof(msg string, fields ...any)
	Errorf(msg string, fields ...any)
}

type API struct {
	storage  storage.Storage
	logger   Logger
	syncSave bool
}

func (a *API) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", a.updateMetric)
	r.Post("/update/", a.updateMetricJSON)
	r.Get("/value/{type}/{name}", a.getMetric)
	r.Post("/value/", a.getMetricJSON)
	r.Get("/", a.info)

	return r
}

func NewAPI(s storage.Storage, log Logger, syncSave bool) *API {
	return &API{
		storage:  s,
		logger:   log,
		syncSave: syncSave,
	}
}

func (a *API) updateMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	switch mtype {
	case Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			a.logger.Infof("invalid value for gauge metric %q: %s", name, value)
			http.Error(w, errMetricValue.Error(), http.StatusBadRequest)
			return
		}
		a.storage.UpdateGauge(name, v)
	case Counter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			a.logger.Infof("invalid value for counter metric %q: %s", name, value)
			http.Error(w, errMetricValue.Error(), http.StatusBadRequest)
			return
		}
		a.storage.UpdateCounter(name, v)
	default:
		a.logger.Infof(InvalidMetricType, mtype)
		http.Error(w, errMetricType.Error(), http.StatusBadRequest)
		return
	}
}

func (a *API) updateMetricJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(ContentType) != ContentTypeJSON {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var metrics model.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		a.logger.Infof("failed to unmarshal metrics: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metrics.MType {
	case Gauge:
		a.storage.UpdateGauge(metrics.ID, *metrics.Value)
	case Counter:
		a.storage.UpdateCounter(metrics.ID, *metrics.Delta)
	default:
		a.logger.Infof(InvalidMetricType, metrics.MType)
		http.Error(w, errMetricType.Error(), http.StatusBadRequest)
		return
	}

	if a.syncSave {
		_ = a.storage.Save()
	}
}

func (a *API) getMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	switch mtype {
	case Gauge:
		v, ok := a.storage.GetGauge(name)
		if !ok {
			a.logger.Infof("gauge metric %q not found", name)
			http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set(ContentType, "text/plain")

		value := strconv.FormatFloat(v, 'f', -1, 64)
		_, err := w.Write([]byte(value))
		if err != nil {
			a.logger.Infof("error writing response: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case Counter:
		v, ok := a.storage.GetCounter(name)
		if !ok {
			a.logger.Infof("counter metric %q not found", name)
			http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set(ContentType, "text/plain")

		value := strconv.FormatInt(v, 10)
		_, err := w.Write([]byte(value))
		if err != nil {
			a.logger.Infof("error writing response: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		a.logger.Infof(InvalidMetricType, mtype)
		http.Error(w, errMetricType.Error(), http.StatusBadRequest)
	}
}

func (a *API) getMetricJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(ContentType) != ContentTypeJSON {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var metrics model.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		a.logger.Infof("failed to unmarshal metrics: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set(ContentType, ContentTypeJSON)

	switch metrics.MType {
	case Gauge:
		v, ok := a.storage.GetGauge(metrics.ID)
		if !ok {
			a.logger.Infof("%s metric %q not found", metrics.MType, metrics.ID)
			http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
			return
		}

		metrics.Value = &v
		if err := json.NewEncoder(w).Encode(&metrics); err != nil {
			a.logger.Infof("failed to marshal %s metrics: %v", metrics.MType, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case Counter:
		v, ok := a.storage.GetCounter(metrics.ID)
		if !ok {
			a.logger.Infof("%s metric %q not found", metrics.MType, metrics.ID)
			http.Error(w, errMetricNotFound.Error(), http.StatusNotFound)
			return
		}
		metrics.Delta = &v

		if err := json.NewEncoder(w).Encode(&metrics); err != nil {
			a.logger.Infof("failed to marshal %s metrics: %v", metrics.MType, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	default:
		a.logger.Infof(InvalidMetricType, metrics.MType)
		http.Error(w, errMetricType.Error(), http.StatusBadRequest)
	}
}

func (a *API) info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(ContentType, ContentTypeHTML)

	_, _ = fmt.Fprintln(w, "<p>Counter metrics:</p><ul>")
	for name, value := range a.storage.GetAllCounter() {
		_, _ = fmt.Fprintf(w, "<li>%s: %v</li>", name, value)
	}
	_, _ = fmt.Fprintln(w, "</ul><p>Gauge metrics:</p><ul>")
	for name, value := range a.storage.GetAllGauge() {
		_, _ = fmt.Fprintf(w, "<li>%s: %v</li>", name, value)
	}
	_, _ = fmt.Fprintln(w, "</ul>")
}
