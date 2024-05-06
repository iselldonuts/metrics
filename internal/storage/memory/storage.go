package memory

import (
	"fmt"
	"strings"
	"sync"
)

type Storage struct {
	GaugeMut   sync.RWMutex
	GaugeMap   map[string]float64
	CounterMut sync.RWMutex
	CounterMap map[string]int64
}

func NewStorage() *Storage {
	return &Storage{
		GaugeMap:   make(map[string]float64),
		GaugeMut:   sync.RWMutex{},
		CounterMap: make(map[string]int64),
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

func (m *Storage) String() string {
	sb := strings.Builder{}
	sb.WriteString("Storage{\n")
	sb.WriteString("\tGauge:\n")
	for n, v := range m.GaugeMap {
		sb.WriteString(fmt.Sprintf("\t\t%s = %f\n", n, v))
	}
	sb.WriteString("\tCounter:\n")
	for n, v := range m.CounterMap {
		sb.WriteString(fmt.Sprintf("\t\t%s = %d\n", n, v))
	}
	sb.WriteString("}")
	return sb.String()
}
