package main

import (
	"os"

	httpserver "github.com/bazookajoe1/metrics-collector/internal/http-server"
	"github.com/bazookajoe1/metrics-collector/internal/memstorage"
)

func main() {
	// TODO: init storage
	servStorage := &memstorage.InMemoryStorage{}
	servStorage.Init()

	// TODO: init http server
	server := &httpserver.HttpServer{}
	err := server.Init("localhost", "8080", servStorage)
	if err != nil {
		os.Exit(-1)
	}

	// TODO: register handlers
	server.RegisterHandler("/update/", server.MetricHandler)

	// TODO: run server
	server.Run()
}
