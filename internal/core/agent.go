package core

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/api"
	"github.com/iselldonuts/metrics/internal/config/agent"
	"github.com/iselldonuts/metrics/internal/metrics"
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

func (a *Agent) Start() {
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

				res, err := client.R().
					SetHeader(api.ContentType, api.ContentTypeJSON).
					SetBody(body).
					Post(url)

				if err != nil {
					log.Printf("Error updating gauge metric %q: %v", m.Name, err)
					continue
				}

				if res.StatusCode() != http.StatusOK {
					log.Printf("Failure updating metrics %q with status code: %d", m.Name, res.StatusCode())
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

				res, err := client.R().
					SetHeader(api.ContentType, api.ContentTypeJSON).
					SetBody(body).
					Post(url)

				if err != nil {
					log.Printf("Error updating counter metrics %q: %v", m.Name, err)
					continue
				}
				if res.StatusCode() != http.StatusOK {
					log.Printf("Failure updating counter metric %q with status code: %d", m.Name, res.StatusCode())
					continue
				}

				if m.Name == "PollCount" {
					a.poller.ResetCounter()
				}
			}
		}
	}
}
