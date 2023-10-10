package main

import (
	httpserver "github.com/bazookajoe1/metrics-collector/internal/http-server"
	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	serverconfig "github.com/bazookajoe1/metrics-collector/internal/server-config"
	"github.com/bazookajoe1/metrics-collector/internal/storage/memstorage"
)

func main() {
	// TODO: create logger
	logger := zlogger.NewZapLogger()
	// TODO: init storage
	servStorage := memstorage.NewInMemoryStorage()

	config := serverconfig.NewConfig(servStorage, logger)

	// TODO: init http server
	server := httpserver.ServerNew(config)

	// TODO: register handlers
	server.InitRoutes()

	// TODO: run server
	server.Run()
}
