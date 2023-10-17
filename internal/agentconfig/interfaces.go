package agentconfig

import (
	"github.com/bazookajoe1/metrics-collector/internal/logger"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
)

type IConfig interface {
	GetAddress() string
	GetPort() string
	GetPI() time.Duration
	GetRI() time.Duration
	GetCollector() collector.MetricCollector
	GetLogger() logger.ILogger
}
