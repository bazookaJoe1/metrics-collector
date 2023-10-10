package agent_config

import (
	"fmt"
	aap "github.com/bazookajoe1/metrics-collector/internal/agent-arg-parser"
	aep "github.com/bazookajoe1/metrics-collector/internal/agent-env-parser"
	"github.com/bazookajoe1/metrics-collector/internal/collector"
	npv "github.com/bazookajoe1/metrics-collector/internal/net-params-validator"
	"log"
	"time"
)

const emptyString = ""
const invalidDuration time.Duration = 123456789

type IParams interface {
	GetAddr() string
	GetPort() string
	GetPI() time.Duration
	GetRI() time.Duration
}

type Config struct {
	address        string
	port           string
	collector      collector.MetricCollector
	pollInterval   time.Duration
	reportInterval time.Duration
	logger         *log.Logger
}

func NewConfig(collector collector.MetricCollector, logger *log.Logger) *Config {
	c := &Config{
		address:        emptyString,
		port:           emptyString,
		collector:      collector,
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
		logger:         logger,
	}

	clArgsParams := aap.ArgParse()
	envParams := aep.EnvParse()
	err := c.UpdateConfig(clArgsParams, envParams)
	if err != nil {
		panic(err)
	}

	return c

}

func (c *Config) UpdateConfig(p ...IParams) error {
	for _, paramInstance := range p {
		pollInterval := paramInstance.GetPI()
		if pollInterval != invalidDuration {
			c.pollInterval = pollInterval
		}

		reportInterval := paramInstance.GetRI()
		if reportInterval != invalidDuration {
			c.reportInterval = reportInterval
		}

		address := paramInstance.GetAddr()
		err := npv.ValidateIP(address)
		if err == nil {
			c.address = address
		}

		port := paramInstance.GetPort()
		err = npv.ValidatePort(port)
		if err == nil {
			c.port = port
		}
	}

	if c.address != emptyString {
		if c.port != emptyString {
			if c.pollInterval != invalidDuration {
				if c.reportInterval != invalidDuration {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("params problem; addr: %v, port: %v, poll interval: %v, report interval: %v",
		c.address, c.port, c.pollInterval, c.reportInterval)
}

func (c *Config) GetAddress() string {
	return c.address
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetPI() time.Duration {
	return c.pollInterval
}

func (c *Config) GetRI() time.Duration {
	return c.reportInterval
}

func (c *Config) GetCollector() collector.MetricCollector {
	return c.collector
}

func (c *Config) GetLogger() *log.Logger {
	return c.logger
}
