package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const (
	logTracePath bool = false // Whether to log info about …?
)

func init() {
	logFactory := func(dataSetName string, out *os.File, loggingInternals bool) *logrus.Entry {
		logger := logrus.New()
		logger.Level = logrus.TraceLevel
		// logger.Formatter = new(ecslogrus.Formatter)
		logger.ReportCaller = loggingInternals
		logger.Out = out
		return logger.WithFields(logrus.Fields{
			"event.dataset": dataSetName, // It tells Kibana what kind of event is logged here.
			"service.type":  `go-module-diagnostic-proxy`,
		})
	}
	globalLogger = logFactory(`err`, os.Stderr, logTracePath) // Ask Marek Czudowski why the globalLogger has such settings?
	httpLogger = logFactory(`apache.access`, os.Stdout, logTracePath)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		ts := time.Now()
		wWrapper := negroni.NewResponseWriter(w)
		next.ServeHTTP(wWrapper, r)
		respStatus := wWrapper.Status()
		thisEventLogger := httpLogger.WithFields(logrus.Fields{
			"event.duration":            fmt.Sprintf("%s", time.Since(ts)),
			"http.request.method":       r.Method,
			"http.response.body.bytes":  wWrapper.Size(),
			"http.response.status_code": respStatus,
			"url.path":                  path,
		})
		switch {
		case 500 <= respStatus && respStatus <= 599:
			thisEventLogger.Error("We have a bug in the code.")
		case 400 <= respStatus && respStatus <= 499:
			thisEventLogger.Error("The request contained an error.")
		case 300 <= respStatus && respStatus <= 308:
			thisEventLogger.Error("We redirected the request.")
		case 200 <= respStatus && respStatus <= 226:
			thisEventLogger.Info("Successfully handled the request.")
		default:
			thisEventLogger.Debug("I do not know what to say.")
		}
	})
}
