package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func produceRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{module:.+}/@v/list", funcModuleProtocol).Methods(http.MethodGet)
	router.HandleFunc("/{module:.+}/@latest", funcModuleProtocol).Methods(http.MethodGet)
	router.HandleFunc("/{module:.+}/@v/{version}.info", funcModuleProtocol).Methods(http.MethodGet)
	router.HandleFunc("/{module:.+}/@v/{version}.mod", funcModuleProtocol).Methods(http.MethodGet)
	router.HandleFunc("/{module:.+}/@v/{version}.zip", funcModuleProtocol).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(goAway)
	return router
}

func callProxy(path string) ([]byte, error) {
	globalLogger.Debug("callProxy — start quering proxy")
	resp, err := http.Get(fmt.Sprintf("https://proxy.golang.org%s", path))
	globalLogger.Debug("callProxy — stop  quering proxy")
	if err != nil {
		return []byte{}, err
	}
	b, err := io.ReadAll(resp.Body)
	return b, err
}

func funcModuleProtocol(w http.ResponseWriter, r *http.Request) {
	globalLogger.Debug("funcModuleProtocol — start")
	path := r.RequestURI
	b, err := callProxy(path)
	w.Write(b)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			"http.request.uri":    path,
			"http.response.error": err.Error(),
		}).Error("funcForVersions — error")
		return
	}
	globalLogger.WithFields(logrus.Fields{
		"http.request.uri":   r.RequestURI,
		"http.response.body": string(b),
	}).Debug("funcForVersions — success")
	globalLogger.Debug("funcModuleProtocol — stop")
}

func funcForVersions(w http.ResponseWriter, r *http.Request) {
	globalLogger.Debug("funcForVersions — start")
	path := r.RequestURI
	b, err := callProxy(path)
	w.Write(b)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			"http.request.uri":    path,
			"http.response.error": err.Error(),
		}).Error("funcForVersions — error")
		return
	}
	globalLogger.WithFields(logrus.Fields{
		"http.request.uri":   r.RequestURI,
		"http.response.body": string(b),
	}).Debug("funcForVersions — success")
	globalLogger.Debug("funcForVersions — stop")
}

func goAway(w http.ResponseWriter, r *http.Request) {
	globalLogger.WithFields(logrus.Fields{
		"http.request.method": r.Method,
		"http.request.uri":    r.RequestURI,
	}).Error(`Stupid`)
	const goAwayContent = `
		<!DOCTYPE html>
		<html>
			<body>
				<p>Go away!</p>
			</body>
		</html>
    `
	fmt.Fprintf(w, goAwayContent)
}
