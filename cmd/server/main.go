package main

import (
	httpserver "github.com/bazookajoe1/metrics-collector/internal/http-server"
	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	"github.com/bazookajoe1/metrics-collector/internal/serverconfig"
	"github.com/bazookajoe1/metrics-collector/internal/storage/memstorage"
)

func main() {
	logger := zlogger.NewZapLogger()

	config := serverconfig.NewConfig() // read config for server

	servStorage := memstorage.NewInMemoryStorage(config, logger) // storage takes FileSaver from config

	server := httpserver.ServerNew(config, servStorage, logger) // init server with config

	server.InitRoutes()

	server.Run()
}
