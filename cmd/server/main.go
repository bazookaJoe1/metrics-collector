package main

import (
	"log"
	"os"

	httpserver "github.com/bazookajoe1/metrics-collector/internal/http-server"
	"github.com/bazookajoe1/metrics-collector/internal/storages/memstorage"
)

func main() {
	// TODO: create logger
	logger := log.New(os.Stdout, "", log.Flags())
	// TODO: init storage
	servStorage := &memstorage.InMemoryStorage{}
	servStorage.Init()

	// TODO: init http server
	server := httpserver.ServerNew("localhost", "8080", servStorage, logger)

	// TODO: register handlers
	server.InitRoutes()

	// TODO: run server
	server.Run()
}
