package internal

import (
	"fmt"
	"strings"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (m MemStorage) String() string {
	sb := strings.Builder{}
	sb.WriteString("MemStorage{\n")
	sb.WriteString("\tGauge:\n")
	for k, v := range m.Gauge {
		sb.WriteString(fmt.Sprintf("\t\t%s = %f\n", k, v))
	}
	sb.WriteString("\tCounter:\n")
	for k, v := range m.Counter {
		sb.WriteString(fmt.Sprintf("\t\t%s = %d\n", k, v))
	}
	sb.WriteString("}")
	return sb.String()
}
