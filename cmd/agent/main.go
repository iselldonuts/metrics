package main

import (
	"fmt"
	"github.com/iselldonuts/metrics/internal/metrics"
	"net/http"
	"time"
)

func main() {
	parseFlags()
	run()
}

func run() {
	col := metrics.NewCollector()

	go func() {
		counter := 0
		for {
			counter += 1
			fmt.Printf("update #%d\n", counter)
			col.Update()
			time.Sleep(time.Duration(options.pollInterval) * time.Second)
		}
	}()

	go func() {
		counter := 0
		for {
			counter += 1
			fmt.Println()
			fmt.Printf("counter: %d\n", counter)

			gm, cm := col.GetAll()
			for _, m := range gm {
				m := m
				go func() {
					url := fmt.Sprintf("http://%s/update/gauge/%s/%f", options.baseURL, m.Name, m.Value)
					res, err := http.Post(url, "text/plain", nil)
					if err != nil {
						return
					}

					defer func() {
						_ = res.Body.Close()
					}()
				}()
			}

			for _, m := range cm {
				m := m
				go func() {
					url := fmt.Sprintf("http://%s/update/counter/%s/%d", options.baseURL, m.Name, m.Value)
					res, err := http.Post(url, "text/plain", nil)
					if err != nil {
						return
					}

					defer func() {
						_ = res.Body.Close()
					}()
				}()
			}

			time.Sleep(time.Duration(options.reportInterval) * time.Second)
		}
	}()

	time.Sleep(time.Duration(1<<63 - 1))
}
