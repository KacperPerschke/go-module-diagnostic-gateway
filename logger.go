package main

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var globalLogger *logrus.Entry

func init() {
	logger := logrus.New()
	logger.Level = logrus.TraceLevel
	logger.Formatter = &prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02T15:04:05.999999Z07:00",
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	logger.ReportCaller = false
	logger.Out = os.Stderr
	globalLogger = logger.WithFields(logrus.Fields{})
}
