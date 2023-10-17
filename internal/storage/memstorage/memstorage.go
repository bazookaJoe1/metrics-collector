package memstorage

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/logger"
	"github.com/bazookajoe1/metrics-collector/internal/serverconfig"
	"github.com/bazookajoe1/metrics-collector/internal/storage/filesaver"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/bazookajoe1/metrics-collector/internal/metric"
)

type InMemoryStorage struct {
	gauge     map[string]string
	counter   map[string]string
	fileSaver filesaver.FileSaver
	logger    logger.ILogger
	mu        sync.RWMutex
}

func NewInMemoryStorage(c serverconfig.IConfig, l logger.ILogger) *InMemoryStorage {
	s := &InMemoryStorage{logger: l, fileSaver: c.GetFileSaver()}
	s.gauge = make(map[string]string)
	s.counter = make(map[string]string)

	s.loadStorageFromFile()

	go s.RunFileSaver()

	return s
}

func (s *InMemoryStorage) ReadMetricValue(mType string, mName string) (string, error) {
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

func (s *InMemoryStorage) ReadAllMetrics() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := ""
	keys := make([]string, 0, len(s.gauge))
	for key := range s.gauge {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		out += fmt.Sprintf("%s: %s\n", key, s.gauge[key])
	}
	keys = make([]string, 0, len(s.gauge))
	for key := range s.counter {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		out += fmt.Sprintf("%s: %s\n", key, s.counter[key])
	}

	return out
}

func (s *InMemoryStorage) ReadAllMetricsJSON() []byte {
	var out bytes.Buffer

	for name := range s.gauge {
		m, err := s.ReadEntireMetric(metric.Gauge, name)
		if err != nil {
			// TODO: logging
			continue
		}
		mC, err := metric.MConnect(m)
		if err != nil {
			// TODO: logging
			continue
		}
		jsonData, err := mC.MarshalJSON()
		if err != nil {
			// TODO: logging
			continue
		}

		out.Write(jsonData)
		out.Write([]byte{'\n'})
	}

	for name := range s.counter {
		m, err := s.ReadEntireMetric(metric.Counter, name)
		if err != nil {
			// TODO: logging
			continue
		}
		mC, err := metric.MConnect(m)
		if err != nil {
			// TODO: logging
			continue
		}
		jsonData, err := mC.MarshalJSON()
		if err != nil {
			// TODO: logging
			continue
		}

		out.Write(jsonData)
		out.Write([]byte{'\n'})
	}

	return out.Bytes()
}

func (s *InMemoryStorage) UpdateMetric(m *metric.Metric) {
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

	if s.fileSaver.SynchronizedSaving { // save to file when storage is updated (synchronized)
		err := s.fileSaver.Save(s.ReadAllMetricsJSON())
		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}

func (s *InMemoryStorage) ReadEntireMetric(mType string, mName string) (*metric.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch mType {
	case metric.Gauge:
		if val, ok := s.gauge[mName]; ok {
			m, err := metric.NewMetric(mName, metric.Gauge, val)
			if err != nil {
				return nil, fmt.Errorf("cannot return metric: %s", mName)
			}
			return m, nil
		}
		return nil, fmt.Errorf("metric %s not found", mType)
	case metric.Counter:
		if val, ok := s.counter[mName]; ok {
			m, err := metric.NewMetric(mName, metric.Counter, val)
			if err != nil {
				return nil, fmt.Errorf("cannot return metric: %s", mName)
			}
			return m, nil
		}
		return nil, fmt.Errorf("metric %s not found", mType)
	}

	return nil, fmt.Errorf("invalid metric type %s", mType)
}

func (s *InMemoryStorage) SetFileSaver(fs filesaver.FileSaver) {
	s.fileSaver = fs
}

// RunFileSaver according to SynchronizedSaving flag starts saving of storage to file with period set in config
func (s *InMemoryStorage) RunFileSaver() {
	if !s.fileSaver.SynchronizedSaving {
		for range s.fileSaver.SaveTicker.C {
			err := s.fileSaver.Save(s.ReadAllMetricsJSON())
			if err != nil {
				s.logger.Error(err.Error())
			} else {
				s.logger.Debug("saved to file")
			}
		}
	}
}

// loadStorageFromFile check whether filesaver.FileSaver Restore flag is set and loads
// storage contents from file provided by filesaver.FileSaver
func (s *InMemoryStorage) loadStorageFromFile() {
	if s.fileSaver.Restore { // is flag set
		cache, err := os.OpenFile(s.fileSaver.FilePath, os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			s.logger.Error(err.Error())
			return
		}
		defer cache.Close()

		scanner := bufio.NewScanner(cache)

		for scanner.Scan() { // read line by line
			data := scanner.Bytes()
			mC := &metric.MConnector{}
			err = mC.UnmarshalJSON(data) // try to unmarshal data into MConnector
			if err != nil {
				s.logger.Error(err.Error())
				continue
			}

			m, err := metric.MDisConnect(mC) // transform
			if err != nil {
				s.logger.Error(err.Error())
				continue
			}

			switch m.GetType() {
			case metric.Gauge:
				s.gauge[m.GetName()] = m.GetValue()
			case metric.Counter:
				s.counter[m.GetName()] = m.GetValue()
			default:
				s.logger.Error(fmt.Sprintf("invalid metric type: %v", m))
				continue
			}

		}

		s.logger.Debug("loaded storage contents from file")

	}
}
