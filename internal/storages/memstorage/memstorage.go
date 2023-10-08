package memstorage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/bazookajoe1/metrics-collector/internal/metric"
)

type inMemoryStorage struct {
	gauge   map[string]string
	counter map[string]string
	mu      sync.RWMutex
}

func NewInMemoryStorage() *inMemoryStorage {
	s := &inMemoryStorage{}
	s.gauge = make(map[string]string)
	s.counter = make(map[string]string)

	return s
}

func (s *inMemoryStorage) ReadMetric(mType string, mName string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch mType {
	case metric.Gauge:
		if val, ok := s.gauge[mName]; ok {
			return val, nil
		}
	case metric.Counter:
		if val, ok := s.counter[mName]; ok {
			return val, nil
		}
	}

	return "", fmt.Errorf("invalid metric type %s", mType)
}

func (s *inMemoryStorage) ReadAllMetrics() string {
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

func (s *inMemoryStorage) UpdateMetric(m *metric.Metric) {
	// мы заранее понимаем, что все параметры правильные, поэтому ничего проверять не будем
	mName, mType, mValue := m.GetParams()
	switch mType {
	case metric.Gauge:
		// enter critical section
		s.mu.Lock()
		s.gauge[mName] = mValue
		s.mu.Unlock()

	case metric.Counter:
		// enter critical section
		s.mu.Lock()
		tempCVal, err := strconv.ParseInt(s.counter[mName], 10, 64)
		if err != nil {
			tempCVal = 0 // если такого ключа еще нет, то вернет ошибку, т.к. строка пустая
		}

		counterIncrement, err := strconv.ParseInt(mValue, 10, 64)
		if err != nil {
			break
		}

		tempCVal += counterIncrement
		s.counter[mName] = strconv.FormatInt(tempCVal, 10)
		s.mu.Unlock()
	}
}
