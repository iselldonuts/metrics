package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUpdateMetric(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "success",
			url:  "/update/counter/x/10",
			want: want{http.StatusOK},
		},
		{
			name: "no metric name",
			url:  "/update/counter/",
			want: want{http.StatusNotFound},
		},
		{
			name: "wrong metric type",
			url:  "/update/wrong/x/10",
			want: want{http.StatusBadRequest},
		},
		{
			name: "wrong metric value",
			url:  "/update/counter/x/wrong",
			want: want{http.StatusBadRequest},
		},
	}

	r := chi.NewRouter()
	l, _ := zap.NewDevelopment()
	log := l.Sugar()
	s := memory.NewStorage(log)
	a := NewAPI(s, log, false)
	r.Mount("/", a.Routes())

	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = http.MethodPost

			req.URL = srv.URL + test.url

			res, err := req.Send()

			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, test.want.code, res.StatusCode())
		})
	}
}
