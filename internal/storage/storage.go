package storage

import (
	"github.com/iselldonuts/metrics/internal/storage/file"
	"github.com/iselldonuts/metrics/internal/storage/memory"
)

type Storage interface {
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllGauge() map[string]float64
	SetAllGauge(gm map[string]float64)
	GetAllCounter() map[string]int64
	SetAllCounter(cm map[string]int64)
	Load() error
	Save() error
}

type Logger interface {
	Infof(msg string, fields ...any)
	Errorf(msg string, fields ...any)
}

func NewStorage(conf Config, log Logger) Storage {
	if conf.Memory != nil {
		return memory.NewStorage(log)
	}
	if conf.File != nil {
		return file.NewStorage(conf.File.Path, log)
	}
	return nil
}
