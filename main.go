package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var globalLogger, httpLogger *logrus.Entry

func main() {
	var listenPort string

	flag.StringVar(&listenPort, "listen-addr", "5000", "server listen address")
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%s", listenPort)
	globalLogger.WithField("addr", addr).Info("starting server")
	if err := http.ListenAndServe(addr, produceRouter()); err != nil {
		logrus.WithField("event", "start server").Fatal(err)
	}
}
