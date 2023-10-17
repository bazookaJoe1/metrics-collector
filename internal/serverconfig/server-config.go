package serverconfig

import (
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/storage/filesaver"
	"sync"
	"time"

	npv "github.com/bazookajoe1/metrics-collector/internal/netparamsvalidator"
	sap "github.com/bazookajoe1/metrics-collector/internal/serverargparser"
	sep "github.com/bazookajoe1/metrics-collector/internal/serverenvparser"
)

const emptyString = ""
const invalidDuration int = 123456789

// IParams is an interface providing methods to fill Config with data
// collected from environment variables and command line arguments
type IParams interface {
	GetAddr() string
	GetPort() string
	GetStoreInterval() int
	GetFilePath() string
	GetRestore() *bool
}

// Config provides configuration data to httpserver.HTTPServer
type Config struct {
	address string // listening IP address
	port    string // listening port
	mu      sync.RWMutex
	fs      filesaver.FileSaver
}

// NewConfig creates Config instance with parameters collected from
// environment variables and command line arguments
func NewConfig() *Config {
	c := &Config{
		address: emptyString,
		port:    emptyString,
		fs:      filesaver.NewFileSaver(300, "./tmp/metrics-db.json", true),
	}

	clArgsParams := sap.ArgParse() // get parameters from command line arguments
	envParams := sep.EnvParse()    // get parameters from environment variables
	err := c.UpdateConfig(clArgsParams, envParams)
	if err != nil {
		panic(err)
	}

	return c

}

// UpdateConfig sequentially fill Config with data provided from
// different sources (cl args and env vars). Last provided data (if exists) overrides
// previous provided data
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

		storeInterval := paramInstance.GetStoreInterval()
		if storeInterval != invalidDuration {
			c.fs.SetStoreInterval(storeInterval)
		}

		filePath := paramInstance.GetFilePath()
		if filePath != emptyString {
			c.fs.SetFilePath(filePath)
		}

		restore := paramInstance.GetRestore()
		if restore != nil {
			c.fs.SetRestore(*restore)
		}

	}

	if c.address != emptyString { // check everything is ok
		if c.port != emptyString {
			if c.fs.StoreInterval != time.Duration(invalidDuration) {
				return nil
			}
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

func (c *Config) GetFileSaver() filesaver.FileSaver {
	return c.fs
}
