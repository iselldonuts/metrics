package storage

import (
	"fmt"
	"github.com/iselldonuts/metrics/internal/storage/memory"
)

type Storage interface {
	fmt.Stringer
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
}

func NewStorage(conf Config) Storage {
	if conf.Memory != nil {
		return memory.NewStorage()
	}
	return nil
}
