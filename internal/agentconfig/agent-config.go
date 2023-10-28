package agentconfig

import (
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/agentargparser"
	"github.com/bazookajoe1/metrics-collector/internal/agentenvparser"
	"github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	"time"
)

const emptyString = ""
const invalidDuration = time.Duration(123456789) * time.Second

type IParams interface {
	GetAddr() string
	GetPort() string
	GetPI() int
	GetRI() int
}

type Config struct {
	address        string
	port           string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewConfig() *Config {
	c := &Config{
		address:        emptyString,
		port:           emptyString,
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
	}

	clArgsParams := agentargparser.ArgParse()
	envParams := agentenvparser.EnvParse()
	err := c.UpdateConfig(clArgsParams, envParams)
	if err != nil {
		panic(err)
	}

	return c

}

func (c *Config) UpdateConfig(p ...IParams) error {
	for _, paramInstance := range p {
		pollInterval := time.Duration(paramInstance.GetPI()) * time.Second
		if pollInterval != invalidDuration {
			c.pollInterval = pollInterval
		}

		reportInterval := time.Duration(paramInstance.GetRI()) * time.Second
		if reportInterval != invalidDuration {
			c.reportInterval = reportInterval
		}

		address := paramInstance.GetAddr()
		err := netparamsvalidator.ValidateIP(address)
		if err == nil {
			c.address = address
		}

		port := paramInstance.GetPort()
		err = netparamsvalidator.ValidatePort(port)
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

func (c *Config) GetPollInterval() time.Duration {
	return c.pollInterval
}

func (c *Config) GetReportInterval() time.Duration {
	return c.reportInterval
}
