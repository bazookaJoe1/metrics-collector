package main

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/agentconfig"
	"github.com/bazookajoe1/metrics-collector/internal/collector"
	"github.com/bazookajoe1/metrics-collector/internal/httpagent"
	"github.com/bazookajoe1/metrics-collector/internal/logging"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
)

// What metrics we need to collect (name, type)
var allowedMetrics = pcstats.Metrics{
	{ID: "Alloc", MType: pcstats.Gauge},
	{ID: "BuckHashSys", MType: pcstats.Gauge},
	{ID: "Frees", MType: pcstats.Gauge},
	{ID: "GCCPUFraction", MType: pcstats.Gauge},
	{ID: "GCSys", MType: pcstats.Gauge},
	{ID: "HeapAlloc", MType: pcstats.Gauge},
	{ID: "HeapIdle", MType: pcstats.Gauge},
	{ID: "HeapInuse", MType: pcstats.Gauge},
	{ID: "HeapObjects", MType: pcstats.Gauge},
	{ID: "HeapReleased", MType: pcstats.Gauge},
	{ID: "HeapSys", MType: pcstats.Gauge},
	{ID: "LastGC", MType: pcstats.Gauge},
	{ID: "Lookups", MType: pcstats.Gauge},
	{ID: "MCacheInuse", MType: pcstats.Gauge},
	{ID: "MCacheSys", MType: pcstats.Gauge},
	{ID: "MSpanInuse", MType: pcstats.Gauge},
	{ID: "MSpanSys", MType: pcstats.Gauge},
	{ID: "Mallocs", MType: pcstats.Gauge},
	{ID: "NextGC", MType: pcstats.Gauge},
	{ID: "NumForcedGC", MType: pcstats.Gauge},
	{ID: "NumGC", MType: pcstats.Gauge},
	{ID: "OtherSys", MType: pcstats.Gauge},
	{ID: "PauseTotalNs", MType: pcstats.Gauge},
	{ID: "StackInuse", MType: pcstats.Gauge},
	{ID: "StackSys", MType: pcstats.Gauge},
	{ID: "Sys", MType: pcstats.Gauge},
	{ID: "TotalAlloc", MType: pcstats.Gauge},
	{ID: "RandomValue", MType: pcstats.Gauge},
	{ID: "Pollcount", MType: pcstats.Counter},
}

func main() {
	mainCtx := context.Background()

	logger := logging.NewZapLogger()

	collectorInst := collector.NewCollector(allowedMetrics, logger)

	config := agentconfig.NewConfig()

	agent := httpagent.AgentNew(config, collectorInst, logger)

	agent.Run(mainCtx)
}
