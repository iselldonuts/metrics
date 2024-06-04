package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/iselldonuts/metrics/internal/api"
	"go.uber.org/zap"
)

func Gzip(zw *gzip.Writer, log *zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get(api.ContentEncoding), "gzip") {
				cr, err := gzip.NewReader(r.Body)
				if err != nil {
					log.Errorf("error creating gzip compress reader: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer func(cr *gzip.Reader) {
					if err := cr.Close(); err != nil {
						log.Errorf("error closing gzip compress reader: %v", err)
					}
				}(cr)
				r.Body = io.NopCloser(cr)
				r.Header.Del("Content-Encoding")
			}

			if !strings.Contains(r.Header.Get(api.AcceptEncoding), "gzip") {
				h.ServeHTTP(w, r)
				return
			}

			zw.Reset(w)
			cw := &compressWriter{
				ResponseWriter: w,
				zw:             zw,
			}
			defer func(zw *gzip.Writer) {
				if err := zw.Close(); err != nil {
					log.Errorf("error closing gzip compress writer: %v", err)
				}
			}(zw)

			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(cw, r)
		})
	}
}

type compressWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p) //nolint:wrapcheck // leads to unexpected behavior
}
