package main

import (
	"fmt"
	"net/http"
	"strings"

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
	router.HandleFunc("/sumdb/sum.golang.org/supported", funcSUMDBWelcome).Methods(http.MethodGet)
	router.HandleFunc("/sumdb/sum.golang.org/lookup/{module_with_version:.+}", funcSUMDBProtocol).Methods(http.MethodGet)
	router.HandleFunc("/sumdb/sum.golang.org/tile/{tile_tail:.+}", funcSUMDBProtocol).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(goAway)
	return router
}

func funcModuleProtocol(w http.ResponseWriter, r *http.Request) {
	host := `proxy.golang.org`
	path := r.RequestURI
	resp, err := callWorld(host, path)
	w.Write(resp.Payload)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			"http.request.uri":    path,
			"http.response.error": err.Error(),
		}).Error(``)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	respContentType := resp.Header.Get("Content-Type")
	if respContentType == "application/zip" {
		resp.LogAsBase64 = true
	}
	globalLogger.WithFields(logrus.Fields{
		"http.request.uri":           r.RequestURI,
		"http.response.body":         resp.toLog(),
		"http.response.status":       resp.StatusCode,
		"http.response.Content-Type": respContentType,
	}).Debug(``)
}

func funcSUMDBWelcome(w http.ResponseWriter, r *http.Request) {
	globalLogger.WithFields(logrus.Fields{
		"http.request.method": r.Method,
		"http.request.uri":    r.RequestURI,
	}).Info(`We do support sumdb.`)
	w.WriteHeader(http.StatusOK)
}

func funcSUMDBProtocol(w http.ResponseWriter, r *http.Request) {
	host := `sum.golang.org`
	path := strings.TrimPrefix(r.RequestURI, `/sumdb/sum.golang.org`)
	resp, err := callWorld(host, path)
	w.Write(resp.Payload)
	if err != nil {
		globalLogger.WithFields(logrus.Fields{
			"http.request.uri":    path,
			"http.response.error": err.Error(),
		}).Error(``)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	if strings.Contains(r.RequestURI, `tile`) {
		resp.LogAsBase64 = true
	}
	globalLogger.WithFields(logrus.Fields{
		"http.request.uri":           r.RequestURI,
		"http.response.body":         resp.toLog(),
		"http.response.status":       resp.StatusCode,
		"http.response.Content-Type": resp.Header.Get("Content-Type"),
	}).Debug(``)
}

func goAway(w http.ResponseWriter, r *http.Request) {
	globalLogger.WithFields(logrus.Fields{
		"http.request.method": r.Method,
		"http.request.uri":    r.RequestURI,
	}).Error(``)
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
