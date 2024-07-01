package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iselldonuts/metrics/internal/storage/memory"
)

const FMode = 0o600

type Logger interface {
	Infof(msg string, fields ...any)
	Errorf(msg string, fields ...any)
}

type Storage struct {
	logger Logger
	path   string
	memory.Storage
}

func NewStorage(path string, log Logger) *Storage {
	return &Storage{
		Storage: *memory.NewStorage(log),
		path:    path,
		logger:  log,
	}
}

func (s *Storage) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		s.logger.Infof("error reading file: %q", s.path)
		return fmt.Errorf("file read error: %w", err)
	}

	metrics := struct {
		GaugeMap   map[string]float64 `json:"gauge"`
		CounterMap map[string]int64   `json:"counter"`
	}{
		GaugeMap:   make(map[string]float64),
		CounterMap: make(map[string]int64),
	}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		s.logger.Errorf("error unmarshalling file storage: %v", err)
		return fmt.Errorf("unmarshall error: %w", err)
	}

	for name, value := range metrics.GaugeMap {
		s.Storage.UpdateGauge(name, value)
	}
	for name, value := range metrics.CounterMap {
		s.Storage.UpdateCounter(name, value)
	}

	return nil
}

func (s *Storage) Save() error {
	metrics := struct {
		GaugeMap   map[string]float64 `json:"gauge"`
		CounterMap map[string]int64   `json:"counter"`
	}{
		GaugeMap:   s.GaugeMap,
		CounterMap: s.CounterMap,
	}

	s.logger.Infof("saving metrics: %v", metrics)

	data, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("marshall error: %w", err)
	}

	if err := os.WriteFile(s.path, data, FMode); err != nil {
		return fmt.Errorf("file save error: %w", err)
	}
	return nil
}
