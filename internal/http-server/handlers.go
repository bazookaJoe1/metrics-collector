package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Metric struct {
	Type  string
	Name  string
	Value string
}

func (serv *_HTTPServer) MetricSave(res http.ResponseWriter, req *http.Request) {

	serv.Logger.Println("Request", req.URL.Path)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metric := Metric{
		Type:  chi.URLParam(req, "type"),
		Name:  chi.URLParam(req, "name"),
		Value: chi.URLParam(req, "value"),
	}

	err := serv.Strg.UpdateMetric(metric.Type, metric.Name, metric.Value)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

	res.Write([]byte{})
}

func (serv *_HTTPServer) MetricRead(res http.ResponseWriter, req *http.Request) {
	var err error
	serv.Logger.Println("Request", req.URL.Path)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metric := Metric{
		Type:  chi.URLParam(req, "type"),
		Name:  chi.URLParam(req, "name"),
		Value: "",
	}

	metric.Value, err = serv.Strg.ReadMetric(metric.Type, metric.Name)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
	}

	res.Write([]byte(metric.Value))

	_ = err
}

func (serv *_HTTPServer) MetricAll(res http.ResponseWriter, req *http.Request) {
	serv.Logger.Println("Request", req.URL.Path)

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metrics := serv.Strg.ReadAllMetrics()
	res.Write([]byte(metrics))
}
