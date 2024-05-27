package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/iselldonuts/metrics/internal/api"
	"go.uber.org/zap"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p) //nolint:wrapcheck // leads to unexpected behavior
}

func (c *compressWriter) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close() //nolint:wrapcheck // leads to unexpected behavior
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err //nolint:wrapcheck // leads to unexpected behavior
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p) //nolint:wrapcheck // leads to unexpected behavior
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err //nolint:wrapcheck // leads to unexpected behavior
	}
	return c.zr.Close() //nolint:wrapcheck // leads to unexpected behavior
}

func Gzip(log *zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ow := w

			// decompress
			if strings.Contains(r.Header.Get(api.ContentEncoding), "gzip") {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					log.Infof("error creating gzip compress reader: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer func() {
					if err := cr.Close(); err != nil {
						log.Infof("error closing gzip compress reader: %v", err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}()
			}

			// compress
			if strings.Contains(r.Header.Get(api.AcceptEncoding), "gzip") {
				cw := newCompressWriter(w)
				ow = cw
				w.Header().Set("Content-Encoding", "gzip")
				defer func() {
					if err := cw.Close(); err != nil {
						log.Infof("error closing gzip compress writer: %v", err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}()
			}

			h.ServeHTTP(ow, r)
		})
	}
}
