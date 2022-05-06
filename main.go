package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var listenPort string

	flag.StringVar(&listenPort, "listen-addr", "5000", "server listen address")
	flag.Parse()

	addr := fmt.Sprintf(":%s", listenPort)
	globalLogger.WithField("addr", addr).Info("starting server")
	if err := http.ListenAndServe(addr, produceRouter()); err != nil {
		globalLogger.WithField("event", "start server").Fatal(err)
	}
}
