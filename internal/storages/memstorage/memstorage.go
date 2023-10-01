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
	gauge   map[string]string
	counter map[string]string
	mu      sync.RWMutex
}

func (s *InMemoryStorage) Init() {
	s.gauge = make(map[string]string)
	s.counter = make(map[string]string)
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
			return val, nil
		}
	case "counter":
		if val, ok := s.counter[mName]; ok {
			return val, nil
		}
	}

	return "", fmt.Errorf("invalid metric type %s", mType)

}

func (s *InMemoryStorage) ReadAllMetrics() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := ""
	for key, val := range s.gauge {
		out += fmt.Sprintf("%s: %s\n", key, val)
	}
	for key, val := range s.counter {
		out += fmt.Sprintf("%s: %s\n", key, val)
	}

	return out
}

func (s *InMemoryStorage) UpdateMetric(mType string, mName string, mValue string) error {
	if !checkMetricName(mName) {
		return fmt.Errorf("invalid metric name %s", mName)
	}

	switch mType {
	case "gauge":
		err := checkGaugeValue(mValue)
		if err != nil {
			return err
		}
		// enter critical section
		s.mu.Lock()
		s.gauge[mName] = mValue
		s.mu.Unlock()

	case "counter":
		counterIncrement, err := checkCounterValue(mValue)
		if err != nil {
			return err
		}
		// enter critical section
		s.mu.Lock()
		tempCVal, err := strconv.ParseInt(s.counter[mName], 10, 64)
		if err != nil { // такое вряд ли случится, но мало ли
			return err
		}
		tempCVal += counterIncrement
		s.counter[mName] = strconv.FormatInt(tempCVal, 10)
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
func checkGaugeValue(value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	return nil
}

// Checks counter metric value is correct to convert into int64
func checkCounterValue(value string) (int64, error) {
	counterVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return counterVal, nil
}
