package memory

import (
	"sync"
)

type Logger interface {
	Infof(msg string, fields ...any)
	Errorf(msg string, fields ...any)
}

type Storage struct {
	GaugeMap   map[string]float64 `json:"gauge"`
	CounterMap map[string]int64   `json:"counter"`
	logger     Logger
	gaugeMut   sync.RWMutex
	counterMut sync.RWMutex
}

func NewStorage(log Logger) *Storage {
	return &Storage{
		GaugeMap:   make(map[string]float64),
		CounterMap: make(map[string]int64),
		gaugeMut:   sync.RWMutex{},
		counterMut: sync.RWMutex{},
		logger:     log,
	}
}

func (m *Storage) UpdateCounter(name string, value int64) {
	m.counterMut.Lock()
	defer m.counterMut.Unlock()

	m.CounterMap[name] += value
}

func (m *Storage) UpdateGauge(name string, value float64) {
	m.gaugeMut.Lock()
	defer m.gaugeMut.Unlock()

	m.GaugeMap[name] = value
}

func (m *Storage) GetGauge(name string) (float64, bool) {
	v, ok := m.GaugeMap[name]
	return v, ok
}

func (m *Storage) GetCounter(name string) (int64, bool) {
	v, ok := m.CounterMap[name]
	return v, ok
}

func (m *Storage) GetAllGauge() map[string]float64 {
	return m.GaugeMap
}

func (m *Storage) GetAllCounter() map[string]int64 {
	return m.CounterMap
}

func (m *Storage) SetAllGauge(gm map[string]float64) {
	m.gaugeMut.Lock()
	defer m.gaugeMut.Unlock()

	m.GaugeMap = gm
}

func (m *Storage) SetAllCounter(cm map[string]int64) {
	m.counterMut.Lock()
	defer m.counterMut.Unlock()

	m.CounterMap = cm
}

func (m *Storage) Load() error {
	// noop
	return nil
}

func (m *Storage) Save() error {
	// noop
	return nil
}
