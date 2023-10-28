package httpserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type HTTPServer struct {
	Server  *echo.Echo
	Config  IConfig
	Storage IStorage
	Logger  ILogger
}

func ServerNew(c IConfig, s IStorage, l ILogger) *HTTPServer {
	serv := &HTTPServer{
		Config:  c,
		Server:  echo.New(),
		Storage: s,
		Logger:  l,
	}

	return serv
}

func (s *HTTPServer) InitRoutes() {
	s.Server.Use(Decompressor)
	s.Server.Use(middleware.Logger())

	s.Server.RouteNotFound("/*", func(c echo.Context) error { return c.NoContent(http.StatusNotFound) })

	s.Server.GET("/", s.SendAllMetrics)

	gUpdate := s.Server.Group("/update")
	gUpdate.POST("/:type/:name/:value", s.ReceiveMetricString)
	gUpdate.POST("/", s.ReceiveMetricJSON)

	gValue := s.Server.Group("/value")
	gValue.GET("/:type/:name", s.SendMetricString)
	gValue.POST("/", s.SendMetricJSON)
}

func (s *HTTPServer) Run() {
	aP := fmt.Sprintf("%s:%s", s.Config.GetAddress(), s.Config.GetPort())
	s.Server.HideBanner = true
	if err := s.Server.Start(aP); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.Logger.Fatal("shutting down the server")
	}
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
