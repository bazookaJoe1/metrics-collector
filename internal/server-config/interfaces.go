package server_config

import (
	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	"github.com/bazookajoe1/metrics-collector/internal/storage"
)

type IConfig interface {
	GetAddress() string
	GetPort() string
	GetStorage() storage.Storage
	GetLogger() zlogger.ILogger
}
