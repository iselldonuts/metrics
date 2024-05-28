package memory

import (
	"sync"
)

type Storage struct {
	GaugeMap   map[string]float64
	CounterMap map[string]int64
	GaugeMut   sync.RWMutex
	CounterMut sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		GaugeMap:   make(map[string]float64),
		CounterMap: make(map[string]int64),
		GaugeMut:   sync.RWMutex{},
		CounterMut: sync.RWMutex{},
	}
}

func (m *Storage) UpdateCounter(name string, value int64) {
	m.CounterMut.Lock()
	defer m.CounterMut.Unlock()

	m.CounterMap[name] += value
}

func (m *Storage) UpdateGauge(name string, value float64) {
	m.GaugeMut.Lock()
	defer m.GaugeMut.Unlock()

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
	m.GaugeMut.Lock()
	defer m.GaugeMut.Unlock()

	m.GaugeMap = gm
}

func (m *Storage) SetAllCounter(cm map[string]int64) {
	m.CounterMut.Lock()
	defer m.CounterMut.Unlock()

	m.CounterMap = cm
}
