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
	MemStats    *runtime.MemStats
	PollCount   int64
	RandomValue float64
}

func NewPoller() *Poller {
	return &Poller{
		MemStats: &runtime.MemStats{},
	}
}

func (m *Poller) Update() {
	m.PollCount++
	m.RandomValue = rand.Float64()
	runtime.ReadMemStats(m.MemStats)
}

func (m *Poller) ResetCounter() {
	m.PollCount = 0
}

func (m *Poller) GetAll() ([]GaugeMetric, []CounterMetric) {
	gm := make([]GaugeMetric, 0, GaugeMetricsCount)
	gm = append(gm,
		GaugeMetric{Name: "Alloc", Value: float64(m.MemStats.Alloc)},
		GaugeMetric{Name: "BuckHashSys", Value: float64(m.MemStats.BuckHashSys)},
		GaugeMetric{Name: "Frees", Value: float64(m.MemStats.Frees)},
		GaugeMetric{Name: "GCCPUFraction", Value: m.MemStats.GCCPUFraction},
		GaugeMetric{Name: "GCSys", Value: float64(m.MemStats.GCSys)},
		GaugeMetric{Name: "HeapAlloc", Value: float64(m.MemStats.HeapAlloc)},
		GaugeMetric{Name: "HeapIdle", Value: float64(m.MemStats.HeapIdle)},
		GaugeMetric{Name: "HeapInuse", Value: float64(m.MemStats.HeapInuse)},
		GaugeMetric{Name: "HeapObjects", Value: float64(m.MemStats.HeapObjects)},
		GaugeMetric{Name: "HeapReleased", Value: float64(m.MemStats.HeapReleased)},
		GaugeMetric{Name: "HeapSys", Value: float64(m.MemStats.HeapSys)},
		GaugeMetric{Name: "LastGC", Value: float64(m.MemStats.LastGC)},
		GaugeMetric{Name: "Lookups", Value: float64(m.MemStats.Lookups)},
		GaugeMetric{Name: "MCacheInuse", Value: float64(m.MemStats.MCacheInuse)},
		GaugeMetric{Name: "MCacheSys", Value: float64(m.MemStats.MCacheSys)},
		GaugeMetric{Name: "MSpanInuse", Value: float64(m.MemStats.MSpanInuse)},
		GaugeMetric{Name: "MSpanSys", Value: float64(m.MemStats.MSpanSys)},
		GaugeMetric{Name: "Mallocs", Value: float64(m.MemStats.Mallocs)},
		GaugeMetric{Name: "NextGC", Value: float64(m.MemStats.NextGC)},
		GaugeMetric{Name: "NumForcedGC", Value: float64(m.MemStats.NumForcedGC)},
		GaugeMetric{Name: "NumGC", Value: float64(m.MemStats.NumGC)},
		GaugeMetric{Name: "OtherSys", Value: float64(m.MemStats.OtherSys)},
		GaugeMetric{Name: "PauseTotalNs", Value: float64(m.MemStats.PauseTotalNs)},
		GaugeMetric{Name: "StackInuse", Value: float64(m.MemStats.StackInuse)},
		GaugeMetric{Name: "StackSys", Value: float64(m.MemStats.StackSys)},
		GaugeMetric{Name: "Sys", Value: float64(m.MemStats.Sys)},
		GaugeMetric{Name: "TotalAlloc", Value: float64(m.MemStats.TotalAlloc)},
		GaugeMetric{Name: "RandomValue", Value: m.RandomValue})

	cm := make([]CounterMetric, 0, CounterMetricsCount)
	cm = append(cm, CounterMetric{Name: "PollCount", Value: m.PollCount})

	return gm, cm
}
