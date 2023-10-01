package memstorage

import (
	"fmt"
	"strconv"
	"sync"
)

type Metric struct {
	Type  string
	Name  string
	Value string
}

type InMemoryStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      sync.RWMutex
}

func (s *InMemoryStorage) Init() {
	s.gauge = make(map[string]float64)
	s.counter = make(map[string]int64)
}

func (s *InMemoryStorage) ReadMetric(mType string, mName string) (string, error) {
	if !checkMetricName(mName) {
		return "", fmt.Errorf("invalid metric name %s", mName)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	switch mType {
	case "gauge":
		if val, ok := s.gauge[mName]; ok {
			return strconv.FormatFloat(val, 'f', 6, 64), nil
		}
	case "counter":
		if val, ok := s.counter[mName]; ok {
			return strconv.FormatInt(val, 10), nil
		}
	}

	return "", fmt.Errorf("invalid metric type %s", mType)

}

func (s *InMemoryStorage) ReadAllMetrics() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := ""
	for key, val := range s.gauge {
		out += fmt.Sprintf("%s: %s\n", key, strconv.FormatFloat(val, 'f', 6, 64))
	}
	for key, val := range s.counter {
		out += fmt.Sprintf("%s: %s\n", key, strconv.FormatInt(val, 10))
	}

	return out
}

func (s *InMemoryStorage) UpdateMetric(mType string, mName string, mValue string) error {
	if !checkMetricName(mName) {
		return fmt.Errorf("invalid metric name %s", mName)
	}

	switch mType {
	case "gauge":
		float64Value, err := checkGaugeValue(mValue)
		if err != nil {
			return err
		}
		// enter critical section
		s.mu.Lock()
		s.gauge[mName] = float64Value
		s.mu.Unlock()

	case "counter":
		int64Value, err := checkCounterValue(mValue)
		if err != nil {
			return err
		}
		// enter critical section
		s.mu.Lock()
		s.counter[mName] += int64Value
		s.mu.Unlock()
	default:
		return fmt.Errorf("non-existent metric type")
	}
	return nil
}

// Checks metric name is not empty
func checkMetricName(name string) bool {
	return name != ""
}

// Checks gauge metric value is correct to convert into float64
func checkGaugeValue(value string) (float64, error) {
	float64Value, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return float64Value, nil
}

// Checks counter metric value is correct to convert into int64
func checkCounterValue(value string) (int64, error) {
	int64Value, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return int64Value, nil
}
