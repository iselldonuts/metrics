package storage

import (
	"github.com/iselldonuts/metrics/internal/storage/memory"
)

type Config struct {
	Memory *memory.Config
}
