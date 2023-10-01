package httpagent

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
	"github.com/go-resty/resty/v2"
)

type _HTTPAgent struct {
	Client          *resty.Client
	Address         string
	Port            string
	Collector       MetricCollector
	PollInterval    time.Duration
	ReportIntervall time.Duration
	Logger          *log.Logger
}

type MetricCollector interface {
	CollectMetrics() error
	GetMetrics() []collector.Metric
	Run(time.Duration)
}

func AgentNew(address string, port string, collector MetricCollector, pollInterval time.Duration, reportInterval time.Duration, logger *log.Logger) *_HTTPAgent {
	return &_HTTPAgent{
		Client:          resty.New(),
		Address:         address,
		Port:            port,
		Collector:       collector,
		PollInterval:    pollInterval,
		ReportIntervall: reportInterval,
		Logger:          logger,
	}
}

func (agent *_HTTPAgent) Run() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go agent.Collector.Run(agent.PollInterval)

	wg.Add(1)
	go func() {
		for {
			time.Sleep(agent.ReportIntervall * time.Second)
			metrics := agent.Collector.GetMetrics()
			agent.sendMetrics(metrics)
		}
	}()

	wg.Wait()
}

func (agent *_HTTPAgent) sendMetrics(metrics []collector.Metric) {
	for _, metric := range metrics {
		endpoint := fmt.Sprintf("%s/%s/%s", metric.MType, metric.MName, metric.MValue)
		url := fmt.Sprintf("http://%s:%s/update/%s", agent.Address, agent.Port, endpoint)

		response, err := agent.Client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			agent.Logger.Println(err)
			continue
		}

		agent.Logger.Println(url, response.StatusCode())

	}
}
