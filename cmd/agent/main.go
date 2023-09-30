package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
	httpagent "github.com/bazookajoe1/metrics-collector/internal/http-agent"
)

func main() {
	logger := log.New(os.Stdout, "", log.Flags())

	collector_inst := &collector.Collector{}
	collector_inst.Init(logger)

	agent := &httpagent.HTTPAgent{
		Client:          &http.Client{Timeout: 5 * time.Second},
		Address:         "localhost",
		Port:            "8080",
		Collector:       collector_inst,
		PollInterval:    2,
		ReportIntervall: 10,
		Logger:          logger,
	}

	agent.Run()
}
