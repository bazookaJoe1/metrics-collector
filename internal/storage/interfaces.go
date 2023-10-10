package storage

import "github.com/bazookajoe1/metrics-collector/internal/metric"

type Storage interface {
	UpdateMetric(*metric.Metric)
	ReadMetric(mType string, mName string) (string, error)
	ReadAllMetrics() string
}
