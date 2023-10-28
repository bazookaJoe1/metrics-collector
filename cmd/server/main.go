package main

import (
	"context"
	"github.com/bazookajoe1/metrics-collector/internal/httpserver"
	"github.com/bazookajoe1/metrics-collector/internal/logging"
	"github.com/bazookajoe1/metrics-collector/internal/serverconfig"
	"github.com/bazookajoe1/metrics-collector/internal/storage/memorystorage"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}

	mainCtx := context.Background()

	logger := logging.NewZapLogger()

	config := serverconfig.NewConfig() // read config for server

	servStorage := memorystorage.NewMemoryStorage(logger, config) // storage takes FileSaver from config

	/*--------------------------------------------------FILESAVER BLOCK START--------------------------------------------------*/
	fsCtx, fsCancel := context.WithCancel(mainCtx) // fileSaver logic runs in goroutine; so we start it and wait until return
	wg.Add(1)
	go func() {
		servStorage.RunFileSaver(fsCtx)
		wg.Done()
	}()
	/*---------------------------------------------------FILESAVER BLOCK END---------------------------------------------------*/

	/*-------------------------------------------------RUN SERVER BLOCK START--------------------------------------------------*/
	server := httpserver.ServerNew(config, servStorage, logger) // init server with config
	server.InitRoutes()
	go server.Run()
	/*--------------------------------------------------RUN SERVER BLOCK END---------------------------------------------------*/

	/*-----------------------------------BLOCKING UNTIL `OS.INTERRUPT RECEIVED` BLOCK START------------------------------------*/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	/*-------------------------------------BLOCKING UNTIL OS.INTERRUPT RECEIVED BLOCK END--------------------------------------*/

	/*---------------------------------------------CANCELLING CONTEXTS BLOCK START---------------------------------------------*/
	fsCancel()
	wg.Wait()
	logger.Info("all contexts done; starting to stop server")
	/*----------------------------------------------CANCELLING CONTEXTS BLOCK END----------------------------------------------*/

	/*---------------------------------------------SHUTTING DOWN THE SERVER BLOCK----------------------------------------------*/
	serverCtx, serverCancel := context.WithTimeout(mainCtx, 5*time.Second)
	defer serverCancel()

	if err := server.Stop(serverCtx); err != nil {
		server.Logger.Fatal(err.Error())
	}
}
