package httpagent

import (
	"context"
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/compressing"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"sync"
	"time"
)

const ContentEncodingGZIP = "gzip"

var ReqHeaderParams = map[string]string{
	echo.HeaderContentType:     echo.MIMEApplicationJSONCharsetUTF8,
	echo.HeaderAcceptEncoding:  ContentEncodingGZIP,
	echo.HeaderContentEncoding: ContentEncodingGZIP,
}

type HTTPAgent struct {
	Client    *resty.Client
	Config    IConfig
	Collector ICollector
	Logger    ILogger
}

func AgentNew(c IConfig, collector ICollector, logger ILogger) *HTTPAgent {
	return &HTTPAgent{
		Client:    resty.New(),
		Config:    c,
		Collector: collector,
		Logger:    logger,
	}
}

// Run starts all logic of HTTPAgent in goroutines, it controls these goroutines and stop them on os.Interrupt signal.
func (a *HTTPAgent) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	collectorCtx, collectorCancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() { // closure; start collector
		a.Collector.Run(collectorCtx, a.Config.GetPollInterval())
		wg.Done()
	}()

	senderCtx, senderCancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() { // closure; send metrics logic in this func
		sendTicker := time.NewTicker(a.Config.GetReportInterval())
		for {
			select {
			case <-sendTicker.C:
				metrics := a.Collector.GetMetrics()
				a.SendMetrics(metrics)
			case <-senderCtx.Done():
				wg.Done()
				a.Logger.Debug("sender context cancelled; returning")
				return
			}
		}
	}()

	/*-----------------------------------BLOCKING UNTIL `OS.INTERRUPT RECEIVED` BLOCK START------------------------------------*/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	/*-------------------------------------BLOCKING UNTIL OS.INTERRUPT RECEIVED BLOCK END--------------------------------------*/

	collectorCancel()
	senderCancel()
	wg.Wait()

	a.Logger.Info("all contexts done; collector and agent stopped")

}

// SendMetrics is the client send function. It takes pcstats.Metrics and convert each pcstats.Metric into JSON, then compresses and send over HTTP.
func (a *HTTPAgent) SendMetrics(metrics pcstats.Metrics) {
	for _, metric := range metrics {
		metricJSON, err := metric.MarshalJSON()
		if err != nil {
			a.Logger.Error(err.Error())
			continue
		}

		gzMetricJSON, err := compressing.GZIPCompress(metricJSON)
		if err != nil {
			a.Logger.Error(err.Error())
			continue
		}

		a.Logger.Debug(fmt.Sprintf("request: %s", string(metricJSON)))

		url := fmt.Sprintf("http://%s:%s/update/", a.Config.GetAddress(), a.Config.GetPort())

		response, err := a.Client.R().SetHeaders(
			ReqHeaderParams,
		).SetBody(gzMetricJSON).Post(url)

		if err != nil {
			a.Logger.Error(err.Error())
		}

		a.Logger.Info(fmt.Sprintf("{\"url\":\"%s\",\"response\":%v,\"status\":%d}",
			url,
			string(response.Body()),
			response.StatusCode()),
		)
	}
}
