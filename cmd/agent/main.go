package main

import (
	agent_config "github.com/bazookajoe1/metrics-collector/internal/agent-config"
	"log"
	"os"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
	httpagent "github.com/bazookajoe1/metrics-collector/internal/http-agent"
	"github.com/bazookajoe1/metrics-collector/internal/metric"
)

// What metrics we need to collect (name, type)
var allowedMetrics = [][2]string{
	{"Alloc", metric.Gauge},
	{"BuckHashSys", metric.Gauge},
	{"Frees", metric.Gauge},
	{"GCCPUFraction", metric.Gauge},
	{"GCSys", metric.Gauge},
	{"HeapAlloc", metric.Gauge},
	{"HeapIdle", metric.Gauge},
	{"HeapInuse", metric.Gauge},
	{"HeapObjects", metric.Gauge},
	{"HeapReleased", metric.Gauge},
	{"HeapSys", metric.Gauge},
	{"LastGC", metric.Gauge},
	{"Lookups", metric.Gauge},
	{"MCacheInuse", metric.Gauge},
	{"MCacheSys", metric.Gauge},
	{"MSpanInuse", metric.Gauge},
	{"MSpanSys", metric.Gauge},
	{"Mallocs", metric.Gauge},
	{"NextGC", metric.Gauge},
	{"NumForcedGC", metric.Gauge},
	{"NumGC", metric.Gauge},
	{"OtherSys", metric.Gauge},
	{"PauseTotalNs", metric.Gauge},
	{"StackInuse", metric.Gauge},
	{"StackSys", metric.Gauge},
	{"Sys", metric.Gauge},
	{"TotalAlloc", metric.Gauge},
	{"RandomValue", metric.Gauge},
	{"Pollcount", metric.Counter},
}

func main() {
	logger := log.New(os.Stdout, "", log.Flags())

	collectorInst := collector.NewCollector(logger, allowedMetrics)

	config := agent_config.NewConfig(collectorInst, logger)

	agent := httpagent.AgentNew(config)

	agent.Run()
}
