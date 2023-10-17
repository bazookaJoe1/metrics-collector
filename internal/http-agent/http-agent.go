package httpagent

import (
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/datacompressor"
	"github.com/bazookajoe1/metrics-collector/internal/logger"
	"github.com/labstack/echo/v4"
	"sync"
	"time"

	agentconfig "github.com/bazookajoe1/metrics-collector/internal/agentconfig"
	"github.com/bazookajoe1/metrics-collector/internal/collector"

	"github.com/bazookajoe1/metrics-collector/internal/metric"
	"github.com/go-resty/resty/v2"
)

const ContentEncodingGZIP = "gzip"

var ReqHeaderParams = map[string]string{
	echo.HeaderContentType:     echo.MIMEApplicationJSONCharsetUTF8,
	echo.HeaderAcceptEncoding:  ContentEncodingGZIP,
	echo.HeaderContentEncoding: ContentEncodingGZIP,
}

type _HTTPAgent struct {
	Client         *resty.Client
	Address        string
	Port           string
	Collector      collector.MetricCollector
	PollInterval   time.Duration
	ReportInterval time.Duration
	Logger         logger.ILogger
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

// sendMetrics conducts cycle JSON marshalling and gzip compressing of metric.Metric structure.
// After that it sends processed data to server
func (agent *_HTTPAgent) sendMetrics(metrics []*metric.Metric) {
	for _, m := range metrics {
		mC, err := metric.MConnect(m)
		if err != nil {
			agent.Logger.Info(err.Error())
			continue
		}

		reqData, err := mC.MarshalJSON()
		if err != nil {
			agent.Logger.Error(err.Error())
			continue
		}

		gzReqData, err := datacompressor.GZIPCompress(reqData) // gzip compressing
		if err != nil {
			agent.Logger.Error(err.Error())
			continue
		}

		agent.Logger.Debug(fmt.Sprintf("plain: %s, compressed: %v", string(reqData), gzReqData))

		url := fmt.Sprintf("http://%s:%s/update/", agent.Address, agent.Port)

		response, err := agent.Client.R().SetHeaders(
			ReqHeaderParams,
		).SetBody(gzReqData).Post(url)

		if err != nil {
			agent.Logger.Error(err.Error())
		}

		agent.Logger.Info(fmt.Sprintf("{\"url\":\"%s\",\"response\":%v,\"status\":%d}",
			url,
			response.Body(),
			response.StatusCode()),
		)
	}
}
