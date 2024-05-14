package storage

import (
	"github.com/iselldonuts/metrics/internal/storage/memory"
)

type Storage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
}

func NewStorage(conf Config) Storage {
	if conf.Memory != nil {
		return memory.NewStorage()
	}
	return nil
}
