package main

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/iselldonuts/metrics/internal/metrics"
)

func main() {
	parseFlags()
	run()
}

func run() {
	fmt.Printf(
		"Running agent | url: %s, ReportInterval: %d, PollInterval: %d\n",
		options.baseURL, options.reportInterval, options.pollInterval,
	)

	poller := metrics.NewPoller()
	client := resty.New()

	updater := func() {
		for {
			poller.Update()
			time.Sleep(time.Duration(options.pollInterval) * time.Second)
		}
	}

	sender := func() {
		doPost := func(url string) {
			_, _ = client.R().
				SetHeader("Content-type", "text/plain").
				Post(url)
		}

		for {
			gm, cm := poller.GetAll()
			for _, m := range gm {
				url := fmt.Sprintf("http://%s/update/gauge/%s/%f", options.baseURL, m.Name, m.Value)
				go doPost(url)
			}

			for _, m := range cm {
				url := fmt.Sprintf("http://%s/update/counter/%s/%d", options.baseURL, m.Name, m.Value)
				go doPost(url)
			}

			time.Sleep(time.Duration(options.reportInterval) * time.Second)
		}
	}

	go updater()
	go sender()

	time.Sleep(time.Duration(1<<63 - 1))
}
