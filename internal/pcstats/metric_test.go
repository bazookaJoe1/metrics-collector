package pcstats

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func Test_checkMetricName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good name",
			args: args{
				name: "TestGoodName",
			},
			wantErr: false,
		},
		{
			name: "Bad name",
			args: args{
				name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkMetricName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("checkMetricName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkMetricType(t *testing.T) {
	type args struct {
		typeName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good type Gauge",
			args: args{
				typeName: Gauge,
			},
			wantErr: false,
		},
		{
			name: "Good type Counter",
			args: args{
				typeName: Counter,
			},
			wantErr: false,
		},
		{
			name: "Bad type random1",
			args: args{
				typeName: "bad type",
			},
			wantErr: true,
		},
		{
			name: "Bad type random2",
			args: args{
				typeName: strconv.FormatInt(int64(rand.Int()), 10),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkMetricType(tt.args.typeName); (err != nil) != tt.wantErr {
				t.Errorf("checkMetricType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkMetricValueGauge(t *testing.T) {
	type args struct {
		value *float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good",
			args: args{
				value: new(float64),
			},
			wantErr: false,
		},
		{
			name: "Bad",
			args: args{
				value: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkMetricValueGauge(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("checkMetricValueGauge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkMetricValueCounter(t *testing.T) {
	type args struct {
		value *int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good",
			args: args{
				value: new(int64),
			},
			wantErr: false,
		},
		{
			name: "Bad",
			args: args{
				value: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkMetricValueCounter(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("checkMetricValueCounter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMetric(t *testing.T) {
	type args struct {
		typeName     string
		name         string
		gaugeValue   *float64
		counterValue *int64
	}
	tests := []struct {
		name    string
		args    args
		want    *Metric
		wantErr bool
	}{
		{
			name: "Good metric Gauge",
			args: args{
				typeName:     Gauge,
				name:         "Test Gauge",
				gaugeValue:   new(float64),
				counterValue: nil,
			},
			want: &Metric{
				ID:    "Test Gauge",
				MType: Gauge,
				Delta: nil,
				Value: new(float64),
			},
			wantErr: false,
		},
		{
			name: "Good metric Counter",
			args: args{
				typeName:     Counter,
				name:         "Test Counter",
				gaugeValue:   nil,
				counterValue: new(int64),
			},
			want: &Metric{
				ID:    "Test Counter",
				MType: Counter,
				Delta: new(int64),
				Value: nil,
			},
			wantErr: false,
		},
		{
			name: "Bad metric Counter (nil value)",
			args: args{
				typeName:     Counter,
				name:         "Bad Counter",
				gaugeValue:   new(float64),
				counterValue: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bad metric Gauge (nil value)",
			args: args{
				typeName:     Gauge,
				name:         "Bad Gauge",
				gaugeValue:   nil,
				counterValue: new(int64),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bad metric (empty name)",
			args: args{
				typeName:     Gauge,
				name:         "",
				gaugeValue:   new(float64),
				counterValue: new(int64),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bad metric (invalid type)",
			args: args{
				typeName:     "randomBadType",
				name:         "GoodMetric",
				gaugeValue:   new(float64),
				counterValue: new(int64),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetric(tt.args.typeName, tt.args.name, tt.args.gaugeValue, tt.args.counterValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}
