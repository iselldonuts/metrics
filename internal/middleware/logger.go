package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

type responseData struct {
	status int
	size   int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, &responseData{
		status: http.StatusOK,
		size:   0,
	}}
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.responseData.status = status
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseData.size += size
	//nolint:wrapcheck // leads to unexpected behavior
	return size, err
}

func Logger(log *zap.SugaredLogger, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := newLoggingResponseWriter(w)
		h(lw, r)

		duration := time.Since(start)

		log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"duration", duration,
			"status", lw.responseData.status,
			"size", lw.responseData.size,
		)
	}
}
