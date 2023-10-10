package httpserver

import (
	"fmt"
	"net/http"
	"os"

	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	"github.com/bazookajoe1/metrics-collector/internal/serverconfig"
	"github.com/bazookajoe1/metrics-collector/internal/storage"

	"github.com/go-chi/chi/v5"
)

type _HTTPServer struct {
	Address string
	Port    string
	Router  *chi.Mux
	Strg    storage.Storage
	Logger  zlogger.ILogger
}

func ServerNew(c serverconfig.IConfig) *_HTTPServer {
	address := c.GetAddress()
	port := c.GetPort()
	strg := c.GetStorage()
	logger := c.GetLogger()
	return &_HTTPServer{
		Address: address,
		Port:    port,
		Strg:    strg,
		Logger:  logger,
		Router:  chi.NewRouter(),
	}
}

func (serv *_HTTPServer) InitRoutes() {
	serv.Router.Use(serv.LogMiddleware)

	serv.Router.Get("/", serv.MetricAll)

	serv.Router.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", serv.MetricSave)
	})
	serv.Router.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", serv.MetricRead)
	})
}

func (serv *_HTTPServer) Run() {
	aP := fmt.Sprintf("%s:%s", serv.Address, serv.Port)
	err := http.ListenAndServe(aP, serv.Router)
	if err != nil {
		serv.Logger.Info(err.Error())
		os.Exit(-1)
	}
}
