package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/iselldonuts/metrics/internal/config/server"
	"github.com/iselldonuts/metrics/internal/storage"
)

const FilePermissions = 0o600

type Archiver struct {
	storage       storage.Storage
	path          string
	storeInterval time.Duration
}

type Metrics struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
}

func NewArchiver(s storage.Storage, conf *server.Config) *Archiver {
	return &Archiver{
		storeInterval: time.Duration(conf.StoreInterval) * time.Second,
		storage:       s,
		path:          conf.FileStoragePath,
	}
}

func (b *Archiver) Start() {
	storeSaveTicker := time.NewTicker(b.storeInterval)
	for {
		<-storeSaveTicker.C
		_ = b.Save()
	}
}

func (b *Archiver) Load() error {
	data, err := os.ReadFile(b.path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", b.path, err)
	}

	var m Metrics
	err = json.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("error unmarshalling file %s: %w", b.path, err)
	}

	b.storage.SetAllGauge(m.Gauge)
	b.storage.SetAllCounter(m.Counter)
	return nil
}

func (b *Archiver) Save() error {
	gm := b.storage.GetAllGauge()
	cm := b.storage.GetAllCounter()

	m := Metrics{
		Gauge:   gm,
		Counter: cm,
	}

	data, err := json.Marshal(m)

	if err != nil {
		return fmt.Errorf("could not marshal gauges: %w", err)
	}
	if err := os.WriteFile(b.path, data, FilePermissions); err != nil {
		return fmt.Errorf("could not write to %q: %w", b.path, err)
	}
	return nil
}
