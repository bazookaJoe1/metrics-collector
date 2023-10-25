package memorystorage

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"os"
	"sync"
)

const (
	emptyString = ""
	gauge       = "gauge"
	counter     = "counter"
)

// MemoryStorage describes storage that stores all metrics in RAM.
type MemoryStorage struct {
	gauge     map[string]float64
	counter   map[string]int64
	logger    ILogger
	fileSaver IFileSaver
	mu        sync.RWMutex
}

func NewMemoryStorage(logger ILogger, config IConfig) *MemoryStorage {
	fs := config.GetFileSaver()
	mS := &MemoryStorage{
		gauge:     make(map[string]float64),
		counter:   make(map[string]int64),
		logger:    logger,
		fileSaver: &fs,
		mu:        sync.RWMutex{},
	}

	mS.loadStorageFromFile()

	return mS
}

// Save performs saving of object that implements IMetric interface
// into RAM storage.
func (s *MemoryStorage) Save(metric pcstats.IMetric) error {
	err := checkMetric(metric)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	s.mu.Lock() // protect storage
	defer s.mu.Unlock()

	switch typeName := metric.GetType(); typeName {
	case gauge:
		value, err := metric.GetGauge()
		if err != nil {
			s.logger.Error(err.Error())
			return err
		}
		s.gauge[metric.GetName()] = *value
	case counter:
		value, err := metric.GetCounter()
		if err != nil {
			s.logger.Error(err.Error())
			return err
		}
		s.counter[metric.GetName()] += *value
	default:
		return fmt.Errorf("broken metric")
	}

	if s.fileSaver.GetSynchronizedFlag() { // if synchronized saving is enabled, we should store storage contents to file
		// every time when metric storage is successful
		metrics := s.GetAll()
		metricsJSON := pcstats.MetricSliceToJSON(metrics)
		err := s.fileSaver.Save(metricsJSON)
		if err != nil {
			s.logger.Error(err.Error())
		} else {
			s.logger.Debug("saved to file")
		}
	}

	return nil
}

// Get tries to read metric from RAM storage according to type name and metric name
func (s *MemoryStorage) Get(typeName string, name string) (*pcstats.Metric, error) {
	err := checkMetricName(name)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	s.mu.RLock() // protect storage
	defer s.mu.RUnlock()

	switch typeName {
	case gauge:
		if value, ok := s.gauge[name]; ok {
			metric, err := pcstats.NewMetric(typeName, name, &value, nil)
			if err != nil {
				s.logger.Error(err.Error()) // here our metric should be valid, but who knows what may happen
				return nil, err
			}
			return metric, nil
		}
		return nil, fmt.Errorf("metric not found")
	case counter:
		if value, ok := s.counter[name]; ok {
			metric, err := pcstats.NewMetric(typeName, name, nil, &value)
			if err != nil {
				s.logger.Error(err.Error()) // here our metric should be valid, but who knows what may happen
				return nil, err
			}
			return metric, nil
		}
		return nil, fmt.Errorf("metric not found")
	}

	return nil, fmt.Errorf("invalid metric type: %s", typeName)
}

// GetAll return all metrics from RAM storage in the form of slice
func (s *MemoryStorage) GetAll() pcstats.Metrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	outMetrics := make([]pcstats.Metric, 0, len(s.gauge)+len(s.counter))
	for name, value := range s.gauge {
		metric, err := pcstats.NewMetric(gauge, name, new(float64), nil)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		err = metric.UpdateGauge(value)
		if err != nil {
			s.logger.Error(err.Error())
		}
		outMetrics = append(outMetrics, *metric)

	}

	for name, value := range s.counter {
		metric, err := pcstats.NewMetric(counter, name, nil, new(int64))
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		err = metric.IncrementCounter(value)
		if err != nil {
			s.logger.Error(err.Error())
		}
		outMetrics = append(outMetrics, *metric)
	}

	return outMetrics
}

// Underlying functionality repeats checks from package pcstats, but it's executed
// upon interface IMetric in purposes of storage. So I decided to repeat it.

// checkMetricName checks that name is not empty string
func checkMetricName(name string) error {
	if name == emptyString {
		return fmt.Errorf("ivalid metric name: %s", name)
	}
	return nil
}

// checkMetricType checks validity of type name. Should be gauge or counter
func checkMetricType(typeName string) error {
	switch typeName {
	case counter:
	case gauge:
	default:
		return fmt.Errorf("invalid metric type: %s", typeName)
	}

	return nil
}

// checkMetricValueGauge checks that pointer on value is not nil
func checkMetricValueGauge(value *float64) error {
	if value == nil {
		return fmt.Errorf("metric type is counter and no value provided for counter")
	}

	return nil
}

// checkMetricValueCounter checks that pointer on value is not nil
func checkMetricValueCounter(value *int64) error {
	if value == nil {
		return fmt.Errorf("metric type is counter and no value provided for counter")
	}

	return nil
}

// checkMetric checks that all parameters of metric is valid
func checkMetric(metric pcstats.IMetric) error {
	err := checkMetricName(metric.GetName())
	if err != nil {
		return err
	}

	err = checkMetricType(metric.GetType())
	if err != nil {
		return err
	}

	switch typeName := metric.GetType(); typeName {
	case counter:
		value, _ := metric.GetCounter()
		err = checkMetricValueCounter(value)
		if err != nil {
			return err
		}
	case gauge:
		value, _ := metric.GetGauge()
		err = checkMetricValueGauge(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("metric is broken")
	}

	return nil
}

// loadStorageFromFile check whether filesaver.FileSaver Restore flag is set and loads
// storage contents from file provided by filesaver.FileSaver
func (s *MemoryStorage) loadStorageFromFile() {
	if s.fileSaver.GetRestore() { // is flag set
		cache, err := os.OpenFile(s.fileSaver.GetFilePath(), os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		defer func() {
			err := cache.Close()
			if err != nil {
				s.logger.Error(err.Error())
			}
		}()

		scanner := bufio.NewScanner(cache)

		for scanner.Scan() { // read line by line
			data := scanner.Bytes()
			metric := pcstats.Metric{}
			err = metric.UnmarshalJSON(data) // try to unmarshal data into Metric
			if err != nil {
				s.logger.Error(err.Error())
				continue
			}

			switch metric.GetType() { // choose where to write (gauge or counter)
			case pcstats.Gauge:
				value, err := metric.GetGauge()
				if err != nil {
					s.logger.Error(err.Error())
					continue // if something wrong - continue (we don't want to dereference nil pointer)
				}
				s.gauge[metric.GetName()] = *value
			case pcstats.Counter:
				value, err := metric.GetCounter()
				if err != nil {
					s.logger.Error(err.Error())
					continue // if something wrong - continue (we don't want to dereference nil pointer)
				}
				s.counter[metric.GetName()] = *value
			default:
				s.logger.Error(fmt.Sprintf("invalid metric type: %v", metric))
				continue
			}

		}

		s.logger.Debug("loaded storage contents from file")

	}
}

// RunFileSaver according to SynchronizedSaving flag starts saving of storage to file with period set in config
func (s *MemoryStorage) RunFileSaver(ctx context.Context) {
	if !s.fileSaver.GetSynchronizedFlag() {
		for {
			select {
			case <-s.fileSaver.GetTicker().C:
				metrics := s.GetAll()
				metricsJSON := pcstats.MetricSliceToJSON(metrics)
				err := s.fileSaver.Save(metricsJSON)
				if err != nil {
					s.logger.Error(err.Error())
				} else {
					s.logger.Debug("saved to file")
				}
			case <-ctx.Done():
				s.fileSaver.GetTicker().Stop()
				metrics := s.GetAll()
				metricsJSON := pcstats.MetricSliceToJSON(metrics)
				err := s.fileSaver.Save(metricsJSON)
				if err != nil {
					s.logger.Error(err.Error())
				}

				s.logger.Debug("file saver context canceled; saved to file")
				return
			}
		}
	}
}
