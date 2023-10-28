package httpagent

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"time"
)

// IConfig is the interface providing methods for configure agent
type IConfig interface {
	GetAddress() string
	GetPort() string
	GetPollInterval() time.Duration
	GetReportInterval() time.Duration
}

// ILogger is the interfaces that allows to work with different loggers.
type ILogger interface {
	Info(string)
	Debug(string)
	Error(string)
	Fatal(string)
}

type ICollector interface {
	Run(ctx context.Context, pollInterval time.Duration)
	GetMetrics() pcstats.Metrics
}
