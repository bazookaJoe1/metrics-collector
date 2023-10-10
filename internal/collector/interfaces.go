package collector

import (
	"github.com/bazookajoe1/metrics-collector/internal/metric"
	"time"
)

type MetricCollector interface {
	CollectMetrics() error
	GetMetrics() []*metric.Metric
	Run(time.Duration)
}
