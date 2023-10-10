package agentconfig

import (
	"log"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
)

type IConfig interface {
	GetAddress() string
	GetPort() string
	GetPI() time.Duration
	GetRI() time.Duration
	GetCollector() collector.MetricCollector
	GetLogger() *log.Logger
}
