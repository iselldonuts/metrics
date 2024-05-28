package core

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/api"
	"github.com/iselldonuts/metrics/internal/config/agent"
	"github.com/iselldonuts/metrics/internal/metrics"
	"go.uber.org/zap"
)

type Agent struct {
	poller         *metrics.Poller
	baseURL        string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func NewAgent(poller *metrics.Poller, conf *agent.Config) *Agent {
	return &Agent{
		poller:         poller,
		baseURL:        conf.Address,
		reportInterval: time.Duration(conf.ReportInterval) * time.Second,
		pollInterval:   time.Duration(conf.PollInterval) * time.Second,
	}
}

func (a *Agent) Start(log *zap.SugaredLogger) {
	client := resty.New()

	pollerTicker := time.NewTicker(a.pollInterval)
	senderTicker := time.NewTicker(a.reportInterval)

	for {
		select {
		case <-pollerTicker.C:
			a.poller.Update()
		case <-senderTicker.C:
			gm, cm := a.poller.GetAll()
			for _, m := range gm {
				value := strconv.FormatFloat(m.Value, 'f', -1, 64)
				url := fmt.Sprintf("http://%s/update/", a.baseURL)

				body := map[string]string{
					"type":  "gauge",
					"id":    m.Name,
					"value": value,
				}

				jsonBody, err := json.Marshal(body)
				if err != nil {
					log.Infof("Error marshalling JSON: %v", err)
					continue
				}

				var buf bytes.Buffer
				gz, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
				if err != nil {
					log.Infof("Unsupported compress level: %v", err)
					continue
				}
				if _, err := gz.Write(jsonBody); err != nil {
					log.Infof("Error writing gzipped data: %v", err)
					continue
				}
				if err := gz.Close(); err != nil {
					log.Infof("Error closing gzip writer: %v", err)
					continue
				}

				res, err := client.R().
					SetHeader(api.ContentType, api.ContentTypeJSON).
					SetHeader(api.ContentEncoding, "gzip").
					SetBody(buf.Bytes()).
					Post(url)

				if err != nil {
					log.Infof("Error updating gauge metric %q: %v", m.Name, err)
					continue
				}

				if res.StatusCode() != http.StatusOK {
					log.Infof("Failure updating metrics %q with status code: %d", m.Name, res.StatusCode())
					continue
				}
			}

			for _, m := range cm {
				value := strconv.FormatInt(m.Value, 10)
				url := fmt.Sprintf("http://%s/update/", a.baseURL)

				body := map[string]string{
					"type":  "counter",
					"id":    m.Name,
					"delta": value,
				}

				jsonBody, err := json.Marshal(body)
				if err != nil {
					log.Infof("Error marshalling JSON: %v", err)
					continue
				}

				var buf bytes.Buffer
				gz, err := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
				if err != nil {
					log.Infof("Unsupported compress level: %v", err)
					continue
				}
				if _, err := gz.Write(jsonBody); err != nil {
					log.Infof("Error writing gzipped data: %v", err)
					continue
				}
				if err := gz.Close(); err != nil {
					log.Infof("Error closing gzip writer: %v", err)
					continue
				}

				res, err := client.R().
					SetHeader(api.ContentType, api.ContentTypeJSON).
					SetHeader(api.ContentEncoding, "gzip").
					SetBody(buf.Bytes()).
					Post(url)

				if err != nil {
					log.Infof("Error updating counter metrics %q: %v", m.Name, err)
					continue
				}
				if res.StatusCode() != http.StatusOK {
					log.Infof("Failure updating counter metric %q with status code: %d", m.Name, res.StatusCode())
					continue
				}

				if m.Name == "PollCount" {
					a.poller.ResetCounter()
				}
			}
		}
	}
}
