// Package metric TODO: change GetParams() to Get[Name, Type, Value]()
package metric

import (
	"fmt"
	"strconv"
)

const emptyString = ""
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
				return nil, err
			}
		case Counter:
			err := checkCounterValue(mValue)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("error metric type: %v", mType)
		}
	} else {
		return nil, fmt.Errorf("error metric name: %v", mName)
	}

	return &Metric{mType: mType, mName: mName, mValue: mValue}, nil
}

func (m *Metric) GetName() string {
	return m.mName
}

func (m *Metric) GetType() string {
	return m.mType
}

func (m *Metric) GetValue() string {
	return m.mValue
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
