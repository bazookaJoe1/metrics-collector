package memorystorage

import (
	"github.com/bazookajoe1/metrics-collector/internal/storage/filesaver"
	"time"
)

type IFileSaver interface {
	Load() ([]byte, error)
	Save(data []byte) error
	GetRestore() bool
	GetFilePath() string
	GetSynchronizedFlag() bool
	GetTicker() *time.Ticker
}

// ILogger is the interfaces that allows to work with different loggers.
type ILogger interface {
	Info(string)
	Debug(string)
	Error(string)
}

type IConfig interface {
	GetFileSaver() filesaver.FileSaver
}
