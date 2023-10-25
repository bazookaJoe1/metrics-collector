package pcstats

import (
	"fmt"
	"sort"
	"strconv"
)

const (
	emptyString = ""
	Gauge       = "gauge"
	Counter     = "counter"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// NewMetric creates and return the instance of Metric. This function automatically validates
// input Metric parameters. You can pass both Gauge and Counter value, but only value according
// to type will be written.
func NewMetric(typeName string, name string, gaugeValue *float64, counterValue *int64) (*Metric, error) {
	if err := checkMetricName(name); err != nil {
		return nil, err
	}

	if err := checkMetricType(typeName); err != nil {
		return nil, err
	}

	switch typeName {
	case Gauge:
		if err := checkMetricValueGauge(gaugeValue); err != nil {
			return nil, err
		}

		metric := &Metric{
			ID:    name,
			MType: typeName,
			Delta: nil,        // according to type one will be nil,
			Value: gaugeValue, // other good
		}

		return metric, nil

	case Counter:
		if err := checkMetricValueCounter(counterValue); err != nil {
			return nil, err
		}

		metric := &Metric{
			ID:    name,
			MType: typeName,
			Delta: counterValue, // according to type one will be nil,
			Value: nil,          // other good
		}

		return metric, nil
	}

	return nil, fmt.Errorf("cannot create metric")

}

func NewMetricFromString(typeName string, name string, value string) (*Metric, error) {
	if err := checkMetricName(name); err != nil {
		return nil, err
	}

	if err := checkMetricType(typeName); err != nil {
		return nil, err
	}

	switch typeName {
	case Gauge:
		gaugeValue, err := StoF64(value)
		if err != nil {
			return nil, err
		}

		if err := checkMetricValueGauge(gaugeValue); err != nil {
			return nil, err
		}

		metric := &Metric{
			ID:    name,
			MType: typeName,
			Delta: nil,        // according to type one will be nil,
			Value: gaugeValue, // other good
		}

		return metric, nil

	case Counter:
		counterValue, err := StoI64(value)
		if err != nil {
			return nil, err
		}

		if err := checkMetricValueCounter(counterValue); err != nil {
			return nil, err
		}

		metric := &Metric{
			ID:    name,
			MType: typeName,
			Delta: counterValue, // according to type one will be nil,
			Value: nil,          // other good
		}

		return metric, nil
	}

	return nil, fmt.Errorf("cannot create metric")

}

// UpdateGauge replaces current Value of metric with new.
func (m *Metric) UpdateGauge(value float64) error {
	if m.MType == Counter {
		return fmt.Errorf("cannot update gauge in counter metric")
	}
	if m.Value != nil {
		*m.Value = value
		return nil
	}
	return fmt.Errorf("metric value is nil")
}

// IncrementCounter increments current Delta of metric on value.
func (m *Metric) IncrementCounter(value int64) error {
	if m.MType == Gauge {
		return fmt.Errorf("cannot update gauge in counter metric")
	}
	if m.Delta != nil {
		*m.Delta += value
		return nil
	}
	return fmt.Errorf("metric value is nil")
}

func (m Metric) GetName() string {
	return m.ID
}

func (m Metric) GetType() string {
	return m.MType
}

// GetGauge returns pointer on Gauge value if Metric type is Gauge. If not returns nil and error.
func (m *Metric) GetGauge() (*float64, error) {
	switch m.MType {
	case Gauge:
		return m.Value, nil
	}
	return nil, fmt.Errorf("metric of type %s has no Counter value", m.MType) // here can only be Counter value
}

// GetCounter returns pointer on Counter value if Metric type is Counter. If not returns nil and error.
func (m *Metric) GetCounter() (*int64, error) {
	switch m.MType {
	case Counter:
		return m.Delta, nil
	}
	return nil, fmt.Errorf("metric of type %s has no Gauge value", m.MType) // here can only be Gauge value
}

func (m *Metric) GetStringValue() (string, error) {
	switch m.GetType() {
	case Gauge:
		value, err := m.GetGauge()
		if err != nil {
			return "", err
		}

		stringValue := strconv.FormatFloat(*value, 'f', -1, 64)

		return stringValue, nil
	case Counter:
		value, err := m.GetCounter()
		if err != nil {
			return "", err
		}

		stringValue := strconv.FormatInt(*value, 10)

		return stringValue, nil
	}

	return "", fmt.Errorf("unhandled error")
}

// checkMetricName checks that name is not empty string
func checkMetricName(name string) error {
	if name == emptyString {
		return fmt.Errorf("ivalid metric name: %s", name)
	}
	return nil
}

// checkMetricType checks validity of type name. Should be Gauge or Counter
func checkMetricType(typeName string) error {
	switch typeName {
	case Counter:
	case Gauge:
	default:
		return fmt.Errorf("invalid metric type: %s", typeName)
	}

	return nil
}

// checkMetricValueGauge checks that pointer on value is not nil
func checkMetricValueGauge(value *float64) error {
	if value == nil {
		return fmt.Errorf("metric type is Gauge and no value provided for Gauge")
	}

	return nil
}

// checkMetricValueCounter checks that pointer on value is not nil
func checkMetricValueCounter(value *int64) error {
	if value == nil {
		return fmt.Errorf("metric type is Counter and no value provided for Counter")
	}

	return nil
}

// StoF64 converts string value to the pointer containing float64 representation of string value.
func StoF64(value string) (*float64, error) {
	tempVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}

	returnVal := new(float64) // we need this var on heap, not on stack
	*returnVal = tempVal

	return returnVal, nil
}

// StoI64 converts string value to the pointer containing int64 representation of string value.
func StoI64(value string) (*int64, error) {
	tempVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	returnVal := new(int64) // we need this var on heap, not on stack
	*returnVal = tempVal

	return returnVal, nil
}

// MetricSliceToString performs conversion of Metric slice to formatted string.
// Each string representation of Metric is separated from another by '\n'.
// Metrics in output string is sorted by type and name.
func MetricSliceToString(metrics Metrics) string {
	outString := ""

	sort.Sort(metrics)

	for _, metric := range metrics {
		value, err := metric.GetStringValue()
		if err != nil {
			continue
		}
		outString += fmt.Sprintf("%s: %s", metric.GetName(), value)
		outString += "\n"
	}

	return outString
}

func MetricSliceToJSON(metrics Metrics) []byte {
	var outData = make([]byte, 0)

	sort.Sort(metrics)

	for _, metric := range metrics {
		data, err := metric.MarshalJSON()
		if err != nil {
			continue
		}
		outData = append(outData, data...)
		outData = append(outData, '\n')
	}

	return outData
}
