package main

import (
	"fmt"
	"github.com/iselldonuts/metrics/internal/metrics"
	"net/http"
	"time"
)

const pollInterval = 2
const reportInterval = 10

func main() {
	col := metrics.NewCollector()

	go func() {
		counter := 0
		for {
			counter += 1
			fmt.Printf("update #%d\n", counter)
			col.Update()
			time.Sleep(pollInterval * time.Second)
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
					url := fmt.Sprintf("http://localhost:8080/update/gauge/%s/%f", m.Name, m.Value)
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
					url := fmt.Sprintf("http://localhost:8080/update/counter/%s/%d", m.Name, m.Value)
					res, err := http.Post(url, "text/plain", nil)
					if err != nil {
						return
					}

					defer func() {
						_ = res.Body.Close()
					}()
				}()
			}

			time.Sleep(reportInterval * time.Second)
		}
	}()

	time.Sleep(time.Duration(1<<63 - 1))
}
