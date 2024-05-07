package metrics

import (
	"math/rand"
	"runtime"
)

const GaugeMetricsCount = 28
const CounterMetricsCount = 1

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

type Poller struct {
	PollCount   int64
	RandomValue float64
	MemStats    *runtime.MemStats
}

func NewPoller() *Poller {
	return &Poller{
		MemStats: &runtime.MemStats{},
	}
}

func (m *Poller) Update() {
	m.PollCount += 1
	m.RandomValue = rand.Float64()
	runtime.ReadMemStats(m.MemStats)
}

func (m *Poller) GetAll() ([]GaugeMetric, []CounterMetric) {
	gm := make([]GaugeMetric, 0, GaugeMetricsCount)
	gm = append(gm, GaugeMetric{Name: "Alloc", Value: float64(m.MemStats.Alloc)})
	gm = append(gm, GaugeMetric{Name: "BuckHashSys", Value: float64(m.MemStats.BuckHashSys)})
	gm = append(gm, GaugeMetric{Name: "Frees", Value: float64(m.MemStats.Frees)})
	gm = append(gm, GaugeMetric{Name: "GCCPUFraction", Value: m.MemStats.GCCPUFraction})
	gm = append(gm, GaugeMetric{Name: "GCSys", Value: float64(m.MemStats.GCSys)})
	gm = append(gm, GaugeMetric{Name: "HeapAlloc", Value: float64(m.MemStats.HeapAlloc)})
	gm = append(gm, GaugeMetric{Name: "HeapIdle", Value: float64(m.MemStats.HeapIdle)})
	gm = append(gm, GaugeMetric{Name: "HeapInuse", Value: float64(m.MemStats.HeapInuse)})
	gm = append(gm, GaugeMetric{Name: "HeapObjects", Value: float64(m.MemStats.HeapObjects)})
	gm = append(gm, GaugeMetric{Name: "HeapReleased", Value: float64(m.MemStats.HeapReleased)})
	gm = append(gm, GaugeMetric{Name: "HeapSys", Value: float64(m.MemStats.HeapSys)})
	gm = append(gm, GaugeMetric{Name: "LastGC", Value: float64(m.MemStats.LastGC)})
	gm = append(gm, GaugeMetric{Name: "Lookups", Value: float64(m.MemStats.Lookups)})
	gm = append(gm, GaugeMetric{Name: "MCacheInuse", Value: float64(m.MemStats.MCacheInuse)})
	gm = append(gm, GaugeMetric{Name: "MCacheSys", Value: float64(m.MemStats.MCacheSys)})
	gm = append(gm, GaugeMetric{Name: "MSpanInuse", Value: float64(m.MemStats.MSpanInuse)})
	gm = append(gm, GaugeMetric{Name: "MSpanSys", Value: float64(m.MemStats.MSpanSys)})
	gm = append(gm, GaugeMetric{Name: "Mallocs", Value: float64(m.MemStats.Mallocs)})
	gm = append(gm, GaugeMetric{Name: "NextGC", Value: float64(m.MemStats.NextGC)})
	gm = append(gm, GaugeMetric{Name: "NumForcedGC", Value: float64(m.MemStats.NumForcedGC)})
	gm = append(gm, GaugeMetric{Name: "NumGC", Value: float64(m.MemStats.NumGC)})
	gm = append(gm, GaugeMetric{Name: "OtherSys", Value: float64(m.MemStats.OtherSys)})
	gm = append(gm, GaugeMetric{Name: "PauseTotalNs", Value: float64(m.MemStats.PauseTotalNs)})
	gm = append(gm, GaugeMetric{Name: "StackInuse", Value: float64(m.MemStats.StackInuse)})
	gm = append(gm, GaugeMetric{Name: "StackSys", Value: float64(m.MemStats.StackSys)})
	gm = append(gm, GaugeMetric{Name: "Sys", Value: float64(m.MemStats.Sys)})
	gm = append(gm, GaugeMetric{Name: "TotalAlloc", Value: float64(m.MemStats.TotalAlloc)})
	gm = append(gm, GaugeMetric{Name: "RandomValue", Value: m.RandomValue})

	cm := make([]CounterMetric, 0, CounterMetricsCount)
	cm = append(cm, CounterMetric{Name: "PollCounter", Value: m.PollCount})
	m.PollCount = 0

	return gm, cm
}
