package main

import (
	"log"
	"net/http"
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
	server := &httpserver.HTTPServer{
		Address: "localhost",
		Port:    "8080",
		Router:  http.NewServeMux(),
		Strg:    servStorage,
		Logger:  logger,
	}

	// TODO: register handlers
	server.Router.HandleFunc("/update/", server.MetricHandler)

	// TODO: run server
	server.Run()
}
