package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/api"
)

type Logger interface {
	Infof(msg string, fields ...any)
	Errorf(msg string, fields ...any)
}

type Sender struct {
	client *resty.Client
	zw     *gzip.Writer
	buf    *bytes.Buffer
	logger Logger
	host   string
}

func NewSender(host string, client *resty.Client, logger Logger) *Sender {
	var buf bytes.Buffer

	return &Sender{
		host:   host,
		logger: logger,
		client: client,
		buf:    &buf,
		zw:     gzip.NewWriter(&buf),
	}
}

func (s *Sender) SendMetric(typ, name, value string) bool {
	defer func() {
		s.buf.Reset()
		s.zw.Reset(s.buf)
	}()

	url := fmt.Sprintf("http://%s/update/", s.host)
	body := map[string]string{
		"type": typ,
		"id":   name,
	}
	if typ == "gauge" {
		body["value"] = value
	} else {
		body["delta"] = value
	}

	s.logger.Infof("body: %v", body)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		s.logger.Errorf("Error marshalling JSON: %v", err)
		return false
	}

	if _, err := s.zw.Write(jsonBody); err != nil {
		s.logger.Errorf("Error writing gzipped data: %v", err)
		return false
	}
	if err := s.zw.Close(); err != nil {
		s.logger.Errorf("Error closing gzip writer: %v", err)
		return false
	}

	res, err := s.client.R().
		SetHeader(api.ContentType, api.ContentTypeJSON).
		SetHeader(api.ContentEncoding, "gzip").
		SetBody(s.buf.Bytes()).
		Post(url)
	if err != nil {
		s.logger.Infof("Error updating %s metric %q: %v", typ, name, err)
		return false
	}

	if res.StatusCode() != http.StatusOK {
		s.logger.Infof("Bad status code updating metrics %q: %d", name, res.StatusCode())
		return false
	}

	return true
}
