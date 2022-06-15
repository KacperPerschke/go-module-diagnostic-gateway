package main

import (
	"fmt"
	"net/http"

	pflag "github.com/spf13/pflag"
)

func main() {
	var configFilePath string
	var listenPort string

	pflag.StringVarP(&configFilePath, "config-file", "c", "./config-default.yaml", "path to the file containing configuration")
	pflag.StringVarP(&listenPort, "listen-addr", "p", "5000", "server listen address")
	pflag.Parse()

	globalConfig = produceConfiguration(configFilePath)

	serverAddress := fmt.Sprintf(":%s", listenPort)
	globalLogger.WithField("serverAddress", serverAddress).Info("starting server")
	if err := http.ListenAndServe(serverAddress, produceRouter()); err != nil {
		globalLogger.WithField("event", "stopping server").Fatal(err)
	}
}
