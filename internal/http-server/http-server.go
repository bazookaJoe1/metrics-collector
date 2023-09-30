package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Storage interface {
	Init()
	UpdateGauge(string, string) error
	UpdateCounter(string, string) error
}

type HTTPServer struct {
	Address string
	Port    string
	Router  *http.ServeMux
	Strg    Storage
	Logger  *log.Logger
}

func (serv *HTTPServer) Run() {
	aP := fmt.Sprintf("%s:%s", serv.Address, serv.Port)
	err := http.ListenAndServe(aP, serv.Router)
	if err != nil {
		os.Exit(-1)
	}
}
