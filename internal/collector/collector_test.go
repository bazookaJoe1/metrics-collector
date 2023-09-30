package collector

import (
	"log"
	"strconv"
	"testing"
)

func TestCollector_CollectMetrics(t *testing.T) {
	c := &Collector{}
	c.Init(log.New(nil, "", 0))

	for counter := 1; counter < 10000; counter++ {
		err := c.CollectMetrics()
		c.mux.RLock()
		if err != nil {
			t.Fatalf("Collect metrics returned an error: %v", err)
		}

		if c.stats["RandomValue"][1] == "0" {
			t.Errorf("RandomValue was not updated")
		}

		if c.stats["Pollcount"][1] != strconv.FormatInt(int64(counter), 10) {
			t.Fatalf("Pollcount is invalid, want: %s, got: %s", c.stats["Pollcount"][1], strconv.FormatInt(int64(counter), 10))
		}
		c.mux.RUnlock()
	}
}

// я потом подумаю как изменить этот тест, т.к. доступ к мапе осуществляется в рандомном порядке
// func TestCollector_GetMetrics(t *testing.T) {
// 	type fields struct {
// 		stats  map[string][2]string
// 		Logger *log.Logger
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   []Metric
// 	}{
// 		{
// 			name: "Simple test",
// 			fields: fields{stats: map[string][2]string{
// 				"Alloc":         {"gauge", "10.500000"},
// 				"BuckHashSys":   {"gauge", "20.500000"},
// 				"Frees":         {"gauge", "30.500000"},
// 				"GCCPUFraction": {"gauge", "0.123000"},
// 				"GCSys":         {"gauge", "40.500000"},
// 				"HeapAlloc":     {"gauge", "50.500000"},
// 				"HeapIdle":      {"gauge", "60.500000"},
// 				"HeapInuse":     {"gauge", "70.500000"},
// 				"HeapObjects":   {"gauge", "80.500000"},
// 				"HeapReleased":  {"gauge", "90.500000"},
// 				"HeapSys":       {"gauge", "100.500000"},
// 				"LastGC":        {"gauge", "110.500000"},
// 				"Lookups":       {"gauge", "120.500000"},
// 				"MCacheInuse":   {"gauge", "130.500000"},
// 				"MCacheSys":     {"gauge", "140.500000"},
// 				"MSpanInuse":    {"gauge", "150.500000"},
// 				"MSpanSys":      {"gauge", "160.500000"},
// 				"Mallocs":       {"gauge", "170.500000"},
// 				"NextGC":        {"gauge", "180.500000"},
// 				"NumForcedGC":   {"gauge", "190.500000"},
// 				"NumGC":         {"gauge", "200.500000"},
// 				"OtherSys":      {"gauge", "210.500000"},
// 				"PauseTotalNs":  {"gauge", "220.500000"},
// 				"StackInuse":    {"gauge", "230.500000"},
// 				"StackSys":      {"gauge", "240.500000"},
// 				"Sys":           {"gauge", "250.500000"},
// 				"TotalAlloc":    {"gauge", "260.500000"},
// 				"RandomValue":   {"gauge", "2.718280"},
// 				"Pollcount":     {"counter", "42"},
// 			}},
// 			want: []Metric{
// 				{"gauge", "Alloc", "10.500000"},
// 				{"gauge", "BuckHashSys", "20.500000"},
// 				{"gauge", "Frees", "30.500000"},
// 				{"gauge", "GCCPUFraction", "0.123000"},
// 				{"gauge", "GCSys", "40.500000"},
// 				{"gauge", "HeapAlloc", "50.500000"},
// 				{"gauge", "HeapIdle", "60.500000"},
// 				{"gauge", "HeapInuse", "70.500000"},
// 				{"gauge", "HeapObjects", "80.500000"},
// 				{"gauge", "HeapReleased", "90.500000"},
// 				{"gauge", "HeapSys", "100.500000"},
// 				{"gauge", "LastGC", "110.500000"},
// 				{"gauge", "Lookups", "120.500000"},
// 				{"gauge", "MCacheInuse", "130.500000"},
// 				{"gauge", "MCacheSys", "140.500000"},
// 				{"gauge", "MSpanInuse", "150.500000"},
// 				{"gauge", "MSpanSys", "160.500000"},
// 				{"gauge", "Mallocs", "170.500000"},
// 				{"gauge", "NextGC", "180.500000"},
// 				{"gauge", "NumForcedGC", "190.500000"},
// 				{"gauge", "NumGC", "200.500000"},
// 				{"gauge", "OtherSys", "210.500000"},
// 				{"gauge", "PauseTotalNs", "220.500000"},
// 				{"gauge", "StackInuse", "230.500000"},
// 				{"gauge", "StackSys", "240.500000"},
// 				{"gauge", "Sys", "250.500000"},
// 				{"gauge", "TotalAlloc", "260.500000"},
// 				{"gauge", "RandomValue", "2.718280"},
// 				{"counter", "Pollcount", "42"},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Collector{
// 				stats:  tt.fields.stats,
// 				Logger: tt.fields.Logger,
// 			}
// 			if got := c.GetMetrics(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Collector.GetMetrics() = %v,\nwant %v", got, tt.want)
// 			}
// 		})
// 	}
// }
