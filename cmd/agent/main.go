package main

import (
	"log"
	"os"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
	httpagent "github.com/bazookajoe1/metrics-collector/internal/http-agent"
)

func main() {
	logger := log.New(os.Stdout, "", log.Flags())

	collectorInst := &collector.Collector{}
	collectorInst.Init(logger)

	agent := httpagent.AgentNew("localhost", "8080", collectorInst, 2, 10, logger)

	agent.Run()
}
