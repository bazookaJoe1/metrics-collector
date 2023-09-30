package httpagent

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bazookajoe1/metrics-collector/internal/collector"
)

type HTTPAgent struct {
	Client          *http.Client
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

func (agent *HTTPAgent) Run() {
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

func (agent *HTTPAgent) sendMetrics(metrics []collector.Metric) {
	for _, metric := range metrics {
		endpoint := fmt.Sprintf("%s/%s/%s", metric.MType, metric.MName, metric.MValue)
		url := fmt.Sprintf("http://%s:%s/update/%s", agent.Address, agent.Port, endpoint)

		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			agent.Logger.Println(err)
			continue
		}

		request.Header.Set("Content-Type", "text/plain")

		response, err := agent.Client.Do(request)
		if err != nil {
			agent.Logger.Println(err)
			continue
		}
		defer response.Body.Close()

		agent.Logger.Println(url, response.StatusCode)

	}
}
