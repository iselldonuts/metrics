package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.url, nil)
			w := httptest.NewRecorder()
			UpdateMetric(w, request)
			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
