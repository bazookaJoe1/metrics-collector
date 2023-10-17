package httpserver

import (
	"github.com/bazookajoe1/metrics-collector/internal/datacompressor"
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/bazookajoe1/metrics-collector/internal/metric"
	_ "github.com/labstack/echo/v4"
)

const emptyString = ""

func (serv *HTTPServer) MetricSave(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)

	m, err := metric.NewMetric(
		c.Param("name"),
		c.Param("type"),
		c.Param("value"),
	)

	if err != nil {
		return c.String(http.StatusBadRequest, emptyString)
	}

	serv.Strg.UpdateMetric(m)

	return c.String(http.StatusOK, emptyString)
}

func (serv *HTTPServer) MetricRead(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)

	value, err := serv.Strg.ReadMetricValue(c.Param("type"), c.Param("name"))
	if err != nil {
		return c.String(http.StatusNotFound, emptyString)
	}

	return c.String(http.StatusOK, value)
}

func (serv *HTTPServer) MetricAll(c echo.Context) error {

	metrics := serv.Strg.ReadAllMetrics()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	return c.String(http.StatusOK, metrics)
}

func (serv *HTTPServer) MetricSaveJSON(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	mC := new(metric.MConnector)
	if err := c.Bind(mC); err != nil {
		return err
	}

	m, err := metric.MDisConnect(mC)
	if err != nil {
		serv.Logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, mC)
	}

	serv.Strg.UpdateMetric(m)

	newValue, err := serv.Strg.ReadMetricValue(m.GetType(), m.GetName())
	if err != nil {
		serv.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}
	_ = m.UpdateMetric(newValue) // here we suppose that metric is valid

	mC, _ = metric.MConnect(m) // here we suppose that metric is valid
	reqData, err := mC.MarshalJSON()
	if err != nil {
		serv.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, emptyString)
	}

	compressedData, err := datacompressor.ServerCompress(c, reqData)
	if err != nil {
		serv.Logger.Error(err.Error()) // if error: reqData returns to compressedData unchanged
	}

	return c.JSONBlob(http.StatusOK, compressedData)
}

func (serv *HTTPServer) MetricReadJSON(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	mC := new(metric.MConnector)
	if err := c.Bind(mC); err != nil {
		return err
	}

	m, err := serv.Strg.ReadEntireMetric(mC.MType, mC.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, emptyString)
	}

	mC, _ = metric.MConnect(m) // here we suppose that metric is valid
	reqData, err := mC.MarshalJSON()
	if err != nil {
		serv.Logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, emptyString)
	}

	compressedData, err := datacompressor.ServerCompress(c, reqData)
	if err != nil {
		serv.Logger.Error(err.Error()) // if error: reqData returns to compressedData unchanged
	}

	return c.JSONBlob(http.StatusOK, compressedData)
}
