package httpagent

import (
	"fmt"
	"log"
	"sync"
	"time"

	agentconfig "github.com/bazookajoe1/metrics-collector/internal/agentconfig"
	"github.com/bazookajoe1/metrics-collector/internal/collector"

	"github.com/bazookajoe1/metrics-collector/internal/metric"
	"github.com/go-resty/resty/v2"
)

type _HTTPAgent struct {
	Client         *resty.Client
	Address        string
	Port           string
	Collector      collector.MetricCollector
	PollInterval   time.Duration
	ReportInterval time.Duration
	Logger         *log.Logger
}

func AgentNew(c agentconfig.IConfig) *_HTTPAgent {
	address := c.GetAddress()
	port := c.GetPort()
	collector := c.GetCollector()
	pollInterval := c.GetPI()
	reportInterval := c.GetRI()
	logger := c.GetLogger()

	return &_HTTPAgent{
		Client:         resty.New(),
		Address:        address,
		Port:           port,
		Collector:      collector,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		Logger:         logger,
	}
}

func (agent *_HTTPAgent) Run() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go agent.Collector.Run(agent.PollInterval)

	wg.Add(1)
	go func() {
		for {
			time.Sleep(agent.ReportInterval)
			metrics := agent.Collector.GetMetrics()
			agent.sendMetrics(metrics)
		}
	}()

	wg.Wait()
}

func (agent *_HTTPAgent) sendMetrics(metrics []*metric.Metric) {
	for _, metric := range metrics {
		mName, mType, mValue := metric.GetParams()
		endpoint := fmt.Sprintf("%s/%s/%s", mType, mName, mValue)
		url := fmt.Sprintf("http://%s:%s/update/%s", agent.Address, agent.Port, endpoint)

		response, err := agent.Client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			agent.Logger.Println(err)
			continue
		}

		agent.Logger.Println(url, response.StatusCode())

	}
}
