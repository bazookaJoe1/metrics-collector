package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type Storage interface {
	Init()
	UpdateMetric(mType string, mName string, mValue string) error
	ReadMetric(mType string, mName string) (string, error)
	ReadAllMetrics() string
}

type _HTTPServer struct {
	Address string
	Port    string
	Router  *chi.Mux
	Strg    Storage
	Logger  *log.Logger
}

func ServerNew(address string, port string, storage Storage, logger *log.Logger) *_HTTPServer {
	return &_HTTPServer{
		Address: address,
		Port:    port,
		Strg:    storage,
		Logger:  logger,
		Router:  chi.NewRouter(),
	}
}

func (serv *_HTTPServer) InitRoutes() {

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
		os.Exit(-1)
	}
}
