package metric

import (
	"fmt"
	"strconv"
)

// MConnector is used for correct translating metric.Metric object into JSON format
type MConnector struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// MConnect translates metric.Metric structure object to MConnector structure object
func MConnect(m *Metric) (*MConnector, error) {
	mc := &MConnector{
		ID:    m.mName,
		MType: m.mType,
	}

	switch m.mType {
	case Counter:
		delta := new(int64)
		err := error(nil)
		*delta, err = strconv.ParseInt(m.mValue, 10, 64)
		if err != nil {
			return nil, err
		}
		mc.Delta = delta

	case Gauge:
		value := new(float64)
		err := error(nil)
		*value, err = strconv.ParseFloat(m.mValue, 64)
		if err != nil {
			return nil, err
		}
		mc.Value = value
	}

	return mc, nil
}

// MDisConnect translates MConnector structure object to metric.Metric structure object
func MDisConnect(mC *MConnector) (*Metric, error) {
	m := new(Metric)

	if !checkMetricName(mC.ID) {
		return nil, fmt.Errorf("bad metric name")
	}
	m.mName = mC.ID

	switch mC.MType {
	case Counter:
		m.mType = mC.MType
		value := strconv.FormatInt(*mC.Delta, 10)
		err := checkCounterValue(value)
		if err != nil {
			return nil, fmt.Errorf("bad counter value")
		}
		m.mValue = value

	case Gauge:
		m.mType = mC.MType
		value := strconv.FormatFloat(*mC.Value, 'f', 10, 64)
		err := checkGaugeValue(value)
		if err != nil {
			return nil, fmt.Errorf("bad counter value")
		}
		m.mValue = value

	default:
		return nil, fmt.Errorf("bad metric value")
	}

	return m, nil
}
