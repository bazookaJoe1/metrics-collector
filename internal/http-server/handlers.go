package httpserver

import (
	"net/http"

	"github.com/bazookajoe1/metrics-collector/internal/metric"
	"github.com/go-chi/chi/v5"
)

func (serv *_HTTPServer) MetricSave(res http.ResponseWriter, req *http.Request) {
	var headerWritten bool

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metric, err := metric.NewMetric(chi.URLParam(req, "name"),
		chi.URLParam(req, "type"),
		chi.URLParam(req, "value"))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

	serv.Strg.UpdateMetric(metric)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
	}

	if !headerWritten {
		res.WriteHeader(http.StatusOK) // we do this because of stupid example in yandex practicum
	}
	res.Write([]byte{})
}

func (serv *_HTTPServer) MetricRead(res http.ResponseWriter, req *http.Request) {
	var err error
	var headerWritten bool

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	value, err := serv.Strg.ReadMetric(chi.URLParam(req, "type"), chi.URLParam(req, "name"))
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		headerWritten = true
	}

	if !headerWritten {
		res.WriteHeader(http.StatusOK) // we do this because of stupid example in yandex practicum
	}
	res.Write([]byte(value))

	_ = err
}

func (serv *_HTTPServer) MetricAll(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	metrics := serv.Strg.ReadAllMetrics()
	res.WriteHeader(http.StatusOK) // we do this because of stupid example in yandex practicum
	res.Write([]byte(metrics))
}
