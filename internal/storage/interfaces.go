package storage

import (
	"github.com/bazookajoe1/metrics-collector/internal/metric"
)

type Storage interface {
	UpdateMetric(*metric.Metric)
	ReadMetricValue(mType string, mName string) (string, error)
	ReadAllMetrics() string
	ReadEntireMetric(mType string, mName string) (*metric.Metric, error)
	ReadAllMetricsJSON() []byte
}
