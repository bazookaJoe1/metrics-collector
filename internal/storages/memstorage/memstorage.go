package memstorage

import (
	"errors"
	"strconv"
	"sync"
)

type InMemoryStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      sync.RWMutex
}

func (s *InMemoryStorage) Init() {
	s.gauge = make(map[string]float64)
	s.counter = make(map[string]int64)
}

func (s *InMemoryStorage) UpdateGauge(name string, value string) error {

	if !checkMetricName(name) {
		return errors.New("invalid gauge metric name")
	}

	float64Value, err := checkGaugeValue(value)
	if err != nil {
		return err
	}

	// enter critical section
	s.mu.Lock()
	s.gauge[name] = float64Value
	s.mu.Unlock()

	return nil
}

func (s *InMemoryStorage) UpdateCounter(name string, value string) error {

	if !checkMetricName(name) {
		return errors.New("invalid counter metric name")
	}

	int64Value, err := checkCounterValue(value)
	if err != nil {
		return err
	}

	// enter critical section
	s.mu.Lock()
	if val, exist := s.counter[name]; exist {
		s.counter[name] = val + int64Value
	} else {
		s.counter[name] = int64Value
	}
	s.mu.Unlock()

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
