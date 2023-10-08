package metric

import (
	"fmt"
	"strconv"
)

const Counter = "counter"
const Gauge = "gauge"

type Metric struct {
	mType  string
	mName  string
	mValue string
}

func NewMetric(mName, mType, mValue string) (*Metric, error) {
	if checkMetricName(mName) {
		switch mType {
		case Gauge:
			err := checkGaugeValue(mValue)
			if err != nil {
				return &Metric{}, err
			}
		case Counter:
			err := checkCounterValue(mValue)
			if err != nil {
				return &Metric{}, err
			}
		default:
			return &Metric{}, fmt.Errorf("error metric type: %v", mType)
		}
	} else {
		return &Metric{}, fmt.Errorf("error metric name: %v", mName)
	}

	return &Metric{mType: mType, mName: mName, mValue: mValue}, nil
}

// Returns metric params in order: name, type, value
func (m *Metric) GetParams() (string, string, string) {
	return m.mName, m.mType, m.mValue
}

func (m *Metric) UpdateMetric(value string) error {
	switch m.mType {
	case Gauge:
		err := checkGaugeValue(value)
		if err != nil {
			return err
		}
		m.mValue = value
	case Counter:
		err := checkCounterValue(value)
		if err != nil {
			return err
		}

		counter, err := strconv.ParseInt(m.mValue, 10, 64)
		if err != nil {
			return err
		}
		counterIncrement, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		counter += counterIncrement
		m.mValue = strconv.FormatInt(counter, 10)
	}
	return nil
}

var AllowedMetrics = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"RandomValue",
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
func checkCounterValue(value string) error {
	_, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	return nil
}
