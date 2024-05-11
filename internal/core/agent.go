package core

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/metrics"
)

type Agent struct {
	poller         *metrics.Poller
	baseURL        string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func NewAgent(poller *metrics.Poller, conf Config) *Agent {
	return &Agent{
		poller:         poller,
		baseURL:        conf.BaseURL,
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
				url := fmt.Sprintf("http://%s/update/gauge/%s/%f", a.baseURL, m.Name, m.Value)
				_, _ = client.R().
					SetHeader("Content-type", "text/plain").
					Post(url)
			}

			for _, m := range cm {
				url := fmt.Sprintf("http://%s/update/counter/%s/%d", a.baseURL, m.Name, m.Value)
				res, err := client.R().
					SetHeader("Content-type", "text/plain").
					Post(url)
				if err != nil {
					log.Println(err)
					continue
				}
				if m.Name == "PollCounter" && res.StatusCode() == http.StatusOK {
					a.poller.ResetCounter()
				}
			}
		}
	}
}
