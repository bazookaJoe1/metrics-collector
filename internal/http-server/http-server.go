package httpserver

import (
	"fmt"
	"github.com/bazookajoe1/metrics-collector/internal/datacompressor"
	zlogger "github.com/bazookajoe1/metrics-collector/internal/logger"
	"github.com/bazookajoe1/metrics-collector/internal/serverconfig"
	"github.com/bazookajoe1/metrics-collector/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTPServer struct {
	Address string
	Port    string
	Server  *echo.Echo
	Strg    storage.Storage
	Logger  zlogger.ILogger
}

func ServerNew(c serverconfig.IConfig, s storage.Storage, l zlogger.ILogger) *HTTPServer {
	address := c.GetAddress()
	port := c.GetPort()
	serv := &HTTPServer{
		Address: address,
		Port:    port,
		Server:  echo.New(),
		Strg:    s,
		Logger:  l,
	}

	return serv
}

func (serv *HTTPServer) InitRoutes() {
	serv.Server.Use(middleware.Logger())
	serv.Server.Use(datacompressor.ServerDecompressor)

	serv.Server.GET("/", serv.MetricAll)

	gUpdate := serv.Server.Group("/update")
	gUpdate.POST("/:type/:name/:value", serv.MetricSave)
	gUpdate.POST("/", serv.MetricSaveJSON)

	gValue := serv.Server.Group("/value")
	gValue.GET("/:type/:name", serv.MetricRead)
	gValue.POST("/", serv.MetricReadJSON)
}

func (serv *HTTPServer) Run() {
	aP := fmt.Sprintf("%s:%s", serv.Address, serv.Port)
	serv.Server.HideBanner = true
	serv.Server.Logger.Fatal(serv.Server.Start(aP))
}
