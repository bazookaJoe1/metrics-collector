package collector

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Metric struct {
	MType  string
	MName  string
	MValue string
}

type Collector struct {
	stats  map[string][2]string
	mux    sync.RWMutex
	Logger *log.Logger
}

func (c *Collector) Init(logger *log.Logger) {
	c.Logger = logger
	//	c.stats = make(map[string][2]string)
	c.stats = map[string][2]string{
		"Alloc":         {"gauge", "0"},
		"BuckHashSys":   {"gauge", "0"},
		"Frees":         {"gauge", "0"},
		"GCCPUFraction": {"gauge", "0"},
		"GCSys":         {"gauge", "0"},
		"HeapAlloc":     {"gauge", "0"},
		"HeapIdle":      {"gauge", "0"},
		"HeapInuse":     {"gauge", "0"},
		"HeapObjects":   {"gauge", "0"},
		"HeapReleased":  {"gauge", "0"},
		"HeapSys":       {"gauge", "0"},
		"LastGC":        {"gauge", "0"},
		"Lookups":       {"gauge", "0"},
		"MCacheInuse":   {"gauge", "0"},
		"MCacheSys":     {"gauge", "0"},
		"MSpanInuse":    {"gauge", "0"},
		"MSpanSys":      {"gauge", "0"},
		"Mallocs":       {"gauge", "0"},
		"NextGC":        {"gauge", "0"},
		"NumForcedGC":   {"gauge", "0"},
		"NumGC":         {"gauge", "0"},
		"OtherSys":      {"gauge", "0"},
		"PauseTotalNs":  {"gauge", "0"},
		"StackInuse":    {"gauge", "0"},
		"StackSys":      {"gauge", "0"},
		"Sys":           {"gauge", "0"},
		"TotalAlloc":    {"gauge", "0"},
		"RandomValue":   {"gauge", "0"},
		"Pollcount":     {"counter", "0"},
	}
}

func (c *Collector) CollectMetrics() error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	c.mux.Lock()
	defer c.mux.Unlock()
	reflectedStatValues := reflect.ValueOf(stats)
	for key := range c.stats {
		val := reflectedStatValues.FieldByName(key)
		if val.IsValid() { // смотрим есть такое поле в струкутуре
			c.stats[key] = [2]string{c.stats[key][0], fmt.Sprintf("%v", reflectedStatValues.FieldByName(key))}
		} else if c.stats[key][0] == "counter" { // сделаем обновления сразу для всех counter
			err := c.counterUpdate(key)
			if err != nil {
				c.Logger.Println(err)
			}
		}
	}

	for { // we don't need zero random value
		randomValue := rand.NormFloat64()
		if randomValue != 0 {
			c.stats["RandomValue"] = [2]string{"gauge", strconv.FormatFloat(float64(randomValue), 'f', 6, 64)}
			break
		}
	}

	return nil
}

func (c *Collector) GetMetrics() []Metric {
	metrics := make([]Metric, 0, len(c.stats))
	c.mux.RLock()
	for key, value := range c.stats {
		metric := Metric{
			MType:  value[0],
			MName:  key,
			MValue: value[1],
		}
		metrics = append(metrics, metric)
	}
	c.mux.RUnlock()
	return metrics

}

func (c *Collector) Run(pollInterval time.Duration) {
	for {
		err := c.CollectMetrics()
		if err != nil {
			c.Logger.Println(err)
		}
		time.Sleep(pollInterval * time.Second)
	}
}

func (c *Collector) counterUpdate(name string) error {
	if _, ok := c.stats[name]; ok {
		if c.stats[name][0] == "counter" {
			counter, err := strconv.ParseInt(c.stats[name][1], 10, 64)
			if err != nil {
				return err
			}

			counter++
			c.stats[name] = [2]string{"counter", strconv.FormatInt(counter, 10)}
			return nil
		}
		return fmt.Errorf("key %s is not type of counter", name)
	}
	return fmt.Errorf("key %s doesn't exist", name)
}
