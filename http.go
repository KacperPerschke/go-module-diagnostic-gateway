package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

/*
 * GET $GOPROXY/<module>/@v/list returns a list of known versions of the given module, one per line.
 * router.HandleFunc("/{module:.+}/@v/{version}.info", version).Methods(http.MethodGet)
 *
 * GET $GOPROXY/<module>/@v/<version>.info returns JSON-formatted metadata about that version of the given module.
 *
 * GET $GOPROXY/<module>/@v/<version>.mod returns the go.mod file for that version of the given module.
 * router.HandleFunc("/{module:.+}/@v/{version}.mod", mod).Methods(http.MethodGet)
 *
 * GET $GOPROXY/<module>/@v/<version>.zip returns the zip archive for that version of the given module.
 * router.HandleFunc("/{module:.+}/@v/{version}.zip", archive).Methods(http.MethodGet)
 *
 * GET $GOPROXY/<module>/@latest returns JSON-formatted metadata about the latest known version of the given module in the same format as <module>/@v/<version>.info.
 *
 * The latest version should be the version of the module the go command may
 * use if <module>/@v/list is empty or no listed version is suitable.
 * <module>/@latest is optional and may not be implemented by a module proxy.
 *
 */
func handlerForHealthCheckers(w http.ResponseWriter, r *http.Request) {
	globalLogger.Debug("handlerForHealth — start")
	globalLogger.WithFields(logrus.Fields{
		"http.request.method": r.Method,
		"http.request.uri":    r.RequestURI,
	}).Debug(`Internals`)
	const healthContent = `
		<!DOCTYPE html>
		<html>
			<body>
				<p>Pocałujta w … wójta!</p>
			</body>
		</html>
    `
	fmt.Fprintf(w, healthContent)
	globalLogger.Debug("handlerForHealth — stop")
}

func handlerForPing(w http.ResponseWriter, r *http.Request) {
	globalLogger.Debug("handlerForPing — start")
	fmt.Fprintf(w, `pong`)
	globalLogger.Debug("handlerForPing — stop")
}

func funcForVersions(w http.ResponseWriter, r *http.Request) {
	globalLogger.Debug("funcForVersions — start")
	module := mux.Vars(r)["module"]
	version := mux.Vars(r)["version"]
	globalLogger.WithFields(logrus.Fields{
		"http.request.params.module":  module,
		"http.request.params.version": version,
		"http.request.uri":            r.RequestURI,
	}).Debug("funcForVersions — paramteres")
	globalLogger.Debug("funcForVersions — start quering proxy")
	resp, err := http.Get(fmt.Sprintf("https://proxy.golang.org%s", r.RequestURI))
	globalLogger.Debug("funcForVersions — stop  quering proxy")
	if err != nil {
		globalLogger.Fatal(err)
		fmt.Fprintf(w, ``)
	} else {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			globalLogger.Fatal(err)
		} else {
			w.Write(b)
			// fmt.Fprintf(w, resp)
		}
	}
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

func produceRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			handlerForHealthCheckers(w, r)
		},
	)
	router.HandleFunc(
		"/ping",
		func(w http.ResponseWriter, r *http.Request) {
			handlerForPing(w, r)
		},
	)
	router.HandleFunc("/{module:.+}/@v/{version}.info", funcForVersions).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(goAway)
	return router
}
