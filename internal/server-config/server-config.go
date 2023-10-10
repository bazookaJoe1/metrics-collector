package server_config

import (
	"fmt"
	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	npv "github.com/bazookajoe1/metrics-collector/internal/net-params-validator"
	sap "github.com/bazookajoe1/metrics-collector/internal/server-arg-parser"
	sep "github.com/bazookajoe1/metrics-collector/internal/server-env-parser"
	"github.com/bazookajoe1/metrics-collector/internal/storage"
	"sync"
)

const emptyString = ""

type IParams interface {
	GetAddr() string
	GetPort() string
}

type Config struct {
	address string
	port    string
	storage storage.Storage
	logger  zlogger.ILogger
	mu      sync.RWMutex
}

func NewConfig(storage storage.Storage, logger zlogger.ILogger) *Config {
	c := &Config{
		address: emptyString,
		port:    emptyString,
		storage: storage,
		logger:  logger}

	clArgsParams := sap.ArgParse()
	envParams := sep.EnvParse()
	err := c.UpdateConfig(clArgsParams, envParams)
	if err != nil {
		panic(err)
	}

	return c

}

func (c *Config) UpdateConfig(p ...IParams) error {
	for _, paramInstance := range p {
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
			return nil
		}
	}

	return fmt.Errorf("params problem; addr: %v, port: %v",
		c.address, c.port)
}

func (c *Config) GetAddress() string {
	return c.address
}

func (c *Config) GetPort() string {
	return c.port
}

func (c *Config) GetStorage() storage.Storage {
	return c.storage
}

func (c *Config) GetLogger() zlogger.ILogger {
	return c.logger
}
