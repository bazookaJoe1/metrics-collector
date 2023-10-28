package httpserver

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
)

// IStorage is the interface that provides methods to work with different types of storages.
type IStorage interface {
	Save(metric pcstats.IMetric) error
	Get(typeName string, name string) (*pcstats.Metric, error)
	GetAll() pcstats.Metrics
	RunFileSaver(ctx context.Context)
}

// ILogger is the interfaces that allows to work with different loggers.
type ILogger interface {
	Info(string)
	Debug(string)
	Error(string)
	Fatal(string)
}

// IConfig is an interface providing methods to configure HTTPServer params.
type IConfig interface {
	GetAddress() string
	GetPort() string
}
