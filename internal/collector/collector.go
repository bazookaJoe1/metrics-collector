package collector

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type Collector struct {
	stats  map[string]*pcstats.Metric
	logger ILogger
	mu     sync.RWMutex
}

// NewCollector create the instance of Collector. Pool of collected metrics is assigned in allowedMetrics.
func NewCollector(allowedMetrics pcstats.Metrics, logger ILogger) *Collector {
	c := &Collector{
		stats:  make(map[string]*pcstats.Metric),
		logger: logger,
	}

	for _, template := range allowedMetrics {
		switch template.MType {
		case pcstats.Gauge:
			metric, err := pcstats.NewMetric(
				template.MType,
				template.ID,
				new(float64),
				nil,
			)
			if err != nil {
				c.logger.Fatal(err.Error())
			}

			c.stats[metric.GetName()] = metric
		case pcstats.Counter:
			metric, err := pcstats.NewMetric(
				template.MType,
				template.ID,
				nil,
				new(int64),
			)
			if err != nil {
				c.logger.Fatal(err.Error())
			}

			c.stats[metric.GetName()] = metric
		}

	}

	return c
}

// Collect collects metric values from runtime.MemStats and updates appropriate metric in Collector.
// All counters increments automatically. RandomValue also is updated.
func (c *Collector) Collect() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	c.mu.Lock()
	defer c.mu.Unlock()

	reflectedStatsValues := reflect.ValueOf(stats) // we use reflection to easy find `ID` of metric in runtime.MemStats
	for key := range c.stats {
		value := reflectedStatsValues.FieldByName(key)

		// update metrics from runtime.MemStats
		if value.IsValid() { // does such field exist and can we typecast it to float.
			switch { // runtime.MemStats contains only Float64 and UInt values (that agreed with allowedMetric)
			case value.CanFloat(): // if value is float, we store it without direct typecasting
				err := c.stats[key].UpdateGauge(value.Float())
				if err != nil {
					c.logger.Error(err.Error())
				}
			case value.CanUint(): // if value is uint, we typecast it to float
				err := c.stats[key].UpdateGauge(float64(value.Uint()))
				if err != nil {
					c.logger.Error(err.Error())
				}
			}
		}

		// update counter if type is counter
		if c.stats[key].GetType() == pcstats.Counter { // if type is counter we update it (by default runtime.MemStats
			// doesn't contain counters)
			err := c.stats[key].IncrementCounter(1)
			if err != nil {
				c.logger.Error(err.Error())
			}
		}

		// update RandomValue
		for { // we don't need zero random value
			randomValue := rand.NormFloat64()
			if randomValue != 0 {
				err := c.stats["RandomValue"].UpdateGauge(randomValue)
				if err != nil {
					c.logger.Error(err.Error())
				}
				break
			}
		}
	}
}

// GetMetrics return all metrics from Collector.
func (c *Collector) GetMetrics() pcstats.Metrics {
	var out = make(pcstats.Metrics, 0, len(c.stats))

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, metric := range c.stats {
		out = append(out, *metric)
	}

	return out
}

// Run starts Collector. Metrics are collected with interval = pollInterval.
func (c *Collector) Run(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			c.logger.Debug("collector context cancelling; return")
			return
		case <-ticker.C:
			c.Collect()
		}
	}
}
