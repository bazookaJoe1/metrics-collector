package agent_config

import (
	"github.com/bazookajoe1/metrics-collector/internal/collector"
	"log"
	"time"
)

type IConfig interface {
	GetAddress() string
	GetPort() string
	GetPI() time.Duration
	GetRI() time.Duration
	GetCollector() collector.MetricCollector
	GetLogger() *log.Logger
}
