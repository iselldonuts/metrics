package core

import (
	"strconv"
	"time"

	"github.com/iselldonuts/metrics/internal/config/agent"
	"github.com/iselldonuts/metrics/internal/metrics"
)

type Poller interface {
	Update()
	ResetCounter()
	GetAll() ([]metrics.GaugeMetric, []metrics.CounterMetric)
}

type Sender interface {
	SendMetric(typ, name, value string) bool
}

type Logger interface {
	Infof(msg string, fields ...any)
}

type Agent struct {
	poller         Poller
	sender         Sender
	logger         Logger
	baseURL        string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func NewAgent(poller Poller, sender Sender, conf *agent.Config, log Logger) *Agent {
	return &Agent{
		poller:         poller,
		sender:         sender,
		baseURL:        conf.Address,
		logger:         log,
		reportInterval: time.Duration(conf.ReportInterval) * time.Second,
		pollInterval:   time.Duration(conf.PollInterval) * time.Second,
	}
}

func (a *Agent) Start() {
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
				if ok := a.sender.SendMetric("gauge", m.Name, value); !ok {
					a.logger.Infof("Failed to send gauge metric %q = %s", m.Name, value)
					continue
				}
			}

			for _, m := range cm {
				value := strconv.FormatInt(m.Value, 10)
				if ok := a.sender.SendMetric("counter", m.Name, value); !ok {
					a.logger.Infof("Failed to send counter metric %q = %s", m.Name, value)
					continue
				}

				if m.Name == "PollCount" {
					a.poller.ResetCounter()
				}
			}
		}
	}
}
