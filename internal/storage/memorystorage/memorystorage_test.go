package memorystorage

import (
	"github.com/bazookajoe1/metrics-collector/internal/logging"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/bazookajoe1/metrics-collector/internal/storage/filesaver"
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
			name: "Good type gauge",
			args: args{
				typeName: gauge,
			},
			wantErr: false,
		},
		{
			name: "Good type counter",
			args: args{
				typeName: counter,
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

func Test_checkMetric(t *testing.T) {
	type args struct {
		metric pcstats.IMetric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good gauge metric",
			args: args{
				&pcstats.Metric{
					ID:    "Good gauge",
					MType: gauge,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: false,
		},
		{
			name: "Good counter metric",
			args: args{
				&pcstats.Metric{
					ID:    "Good counter",
					MType: counter,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Bad metric (empty name)",
			args: args{
				&pcstats.Metric{
					ID:    "",
					MType: gauge,
					Delta: new(int64),
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Bad gauge metric",
			args: args{
				&pcstats.Metric{
					ID:    "Test gauge",
					MType: gauge,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "Bad counter metric",
			args: args{
				&pcstats.Metric{
					ID:    "Test counter",
					MType: counter,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Bad type metric",
			args: args{
				&pcstats.Metric{
					ID:    "Test",
					MType: "abracadabra",
					Delta: new(int64),
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Empty metric",
			args: args{
				&pcstats.Metric{
					ID:    "",
					MType: "",
					Delta: nil,
					Value: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkMetric(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("checkMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemoryStorage_Save(t *testing.T) {
	type args struct {
		metric pcstats.IMetric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good gauge",
			args: args{
				&pcstats.Metric{
					ID:    "TestGauge",
					MType: gauge,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: false,
		},
		{
			name: "Good counter",
			args: args{
				&pcstats.Metric{
					ID:    "TestCounter",
					MType: counter,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Bad metric (empty name)",
			args: args{
				&pcstats.Metric{
					ID:    "",
					MType: gauge,
					Delta: new(int64),
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Bad gauge metric",
			args: args{
				&pcstats.Metric{
					ID:    "Test gauge",
					MType: gauge,
					Delta: new(int64),
					Value: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "Bad counter metric",
			args: args{
				&pcstats.Metric{
					ID:    "Test counter",
					MType: counter,
					Delta: nil,
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Bad type metric",
			args: args{
				&pcstats.Metric{
					ID:    "Test",
					MType: "abracadabra",
					Delta: new(int64),
					Value: new(float64),
				},
			},
			wantErr: true,
		},
		{
			name: "Empty metric",
			args: args{
				&pcstats.Metric{
					ID:    "",
					MType: "",
					Delta: nil,
					Value: nil,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStorage(MockLogger{}, new(MockFileSaver))
			if err := s.Save(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type MockFileSaver struct{}

func (fs *MockFileSaver) GetFileSaver() filesaver.FileSaver { return filesaver.FileSaver{} }

type MockLogger struct{}

func (l MockLogger) Info(string)  {}
func (l MockLogger) Error(string) {}
func (l MockLogger) Debug(string) {}

func TestMemoryStorage_Get(t *testing.T) {
	wants := []*pcstats.Metric{
		{
			ID:    "TestGauge",
			MType: gauge,
			Delta: nil,
			Value: new(float64),
		},
		{
			ID:    "TestCounter",
			MType: counter,
			Delta: new(int64),
			Value: nil,
		},
		nil,
	}

	*wants[0].Value = 1.001
	*wants[1].Delta = 10

	type args struct {
		typeName string
		name     string
	}
	tests := []struct {
		name    string
		args    args
		want    *pcstats.Metric
		wantErr bool
	}{
		{
			name: "Good gauge",
			args: args{
				typeName: gauge,
				name:     "TestGauge",
			},
			want:    wants[0],
			wantErr: false,
		},
		{
			name: "Good counter",
			args: args{
				typeName: counter,
				name:     "TestCounter",
			},
			want:    wants[1],
			wantErr: false,
		},
		{
			name: "No such metric",
			args: args{
				typeName: counter,
				name:     "NoSuch",
			},
			want:    wants[2],
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStorage(logging.NewZapLogger(), new(MockFileSaver))
			s.gauge["TestGauge"] = 1.001
			s.counter["TestCounter"] = 10
			got, err := s.Get(tt.args.typeName, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				if tt.wantErr == true { // if we want error no need to compare next
					return
				}
			}
			if !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.MType, tt.want.MType) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			switch tt.args.typeName {
			case gauge:
				if !reflect.DeepEqual(got.Delta, tt.want.Delta) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(*got.Value, *tt.want.Value) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			case counter:
				if !reflect.DeepEqual(got.Value, tt.want.Value) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(*got.Delta, *tt.want.Delta) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
