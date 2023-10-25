package httpserver

import (
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/pcstats"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

const emptyString = ""

// ReceiveMetricString is the handler responsible for receiving metrics from request uri.
func (s *HTTPServer) ReceiveMetricString(c echo.Context) error {
	metric, err := pcstats.NewMetricFromString( // try to create metric from string values that we get from uri
		c.Param("type"),
		c.Param("name"),
		c.Param("value"),
	)
	if err != nil {
		s.Logger.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	err = s.Storage.Save(metric)
	if err != nil {
		s.Logger.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError) // this case should not occur, but who knows
	}

	return c.NoContent(http.StatusOK)
}

// SendMetricString is the handler responsible for sending metric value got from storage with parameters from uri.
func (s *HTTPServer) SendMetricString(c echo.Context) error {
	metric, err := s.Storage.Get(c.Param("type"), c.Param("name"))
	if err != nil {
		s.Logger.Debug(err.Error())
		return c.NoContent(http.StatusNotFound)
	}

	stringValue, err := metric.GetStringValue()
	if err != nil {
		s.Logger.Debug(err.Error())
		return c.NoContent(http.StatusInternalServerError) // this case should not occur, but who knows
	}

	return c.String(http.StatusOK, stringValue)
}

// SendAllMetrics is the handler responsible for sending back all metrics got from storage in string format.
func (s *HTTPServer) SendAllMetrics(c echo.Context) error {
	metrics := s.Storage.GetAll()

	responseString := pcstats.MetricSliceToString(metrics)

	return c.String(http.StatusOK, responseString)
}

// ReceiveMetricJSON is the handler responsible for receiving metrics in JSON format.
func (s *HTTPServer) ReceiveMetricJSON(c echo.Context) error {
	metric := new(pcstats.Metric)

	if err := c.Bind(metric); err != nil {
		return err
	} // binding request body to pcstats.Metric

	err := s.Storage.Save(metric)
	if err != nil {
		s.Logger.Error(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	responseMetric, err := s.Storage.Get(metric.GetType(), metric.GetName())
	if err != nil {
		s.Logger.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	sendData, err := responseMetric.MarshalJSON() // serializing into JSON
	if err != nil {
		s.Logger.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}

	compressedData := Compressor(c, sendData)

	body := c.Request().Body
	if body != nil {
		data, err := io.ReadAll(body)
		if err == nil {
			s.Logger.Debug(fmt.Sprintf("request: %s", data))
		}
	}

	return c.JSONBlob(http.StatusOK, compressedData)
}

// SendMetricJSON is the handler responsible for getting metric from storage according to parameters got
// from body in JSON format.
func (s *HTTPServer) SendMetricJSON(c echo.Context) error {
	metric := new(pcstats.Metric)

	if err := c.Bind(metric); err != nil {
		return err
	}

	responseMetric, err := s.Storage.Get(metric.GetType(), metric.GetName())
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	sendData, err := responseMetric.MarshalJSON()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	compressedData := Compressor(c, sendData)

	s.Logger.Debug(fmt.Sprintf("response: %s", sendData))

	return c.JSONBlob(http.StatusOK, compressedData)
}
