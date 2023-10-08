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

	"github.com/bazookajoe1/metrics-collector/internal/metric"
)

type collector struct {
	stats  map[string]*metric.Metric
	mux    sync.RWMutex
	Logger *log.Logger
}

// Create instance of collector and return it. Specify needed metrics in allowedMetrics in the format: [][2]string{ {name, type}, ... }
func NewCollector(logger *log.Logger, allowedMetrics [][2]string) *collector {
	c := &collector{Logger: logger}
	c.stats = make(map[string]*metric.Metric)
	for _, template := range allowedMetrics {
		metric, err := metric.NewMetric(template[0], template[1], "0")
		if err != nil {
			c.Logger.Fatal(err)
		}
		c.stats[template[0]] = metric
	}

	return c
}

func (c *collector) CollectMetrics() error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	c.mux.Lock()
	defer c.mux.Unlock()
	reflectedStatValues := reflect.ValueOf(stats)
	for key := range c.stats {
		val := reflectedStatValues.FieldByName(key)
		if val.IsValid() { // смотрим есть такое поле в струкутуре
			err := c.stats[key].UpdateMetric(fmt.Sprintf("%v", reflectedStatValues.FieldByName(key)))
			if err != nil {
				c.Logger.Println(err)
			}
		}
		_, mType, _ := c.stats[key].GetParams()
		if mType == metric.Counter { // сделаем обновления сразу для всех counter
			err := c.stats[key].UpdateMetric("1")
			if err != nil {
				c.Logger.Println(err)
			}
		}
	}

	for { // we don't need zero random value
		randomValue := rand.NormFloat64()
		if randomValue != 0 {
			c.stats["RandomValue"].UpdateMetric(strconv.FormatFloat(float64(randomValue), 'f', 3, 64))
			break
		}
	}

	return nil
}

func (c *collector) GetMetrics() []*metric.Metric {
	metrics := make([]*metric.Metric, 0, len(c.stats))
	c.mux.RLock()
	for _, metric := range c.stats {
		metrics = append(metrics, metric)
	}
	c.mux.RUnlock()
	return metrics

}

func (c *collector) Run(pollInterval time.Duration) {
	for {
		err := c.CollectMetrics()
		if err != nil {
			c.Logger.Println(err)
		}
		time.Sleep(pollInterval * time.Second)
	}
}
