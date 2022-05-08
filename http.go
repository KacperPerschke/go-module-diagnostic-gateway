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
	router.HandleFunc("/health", handlerForHealthCheckers).Methods(http.MethodGet)
	router.HandleFunc("/ping", handlerForPing).Methods(http.MethodGet)
	router.HandleFunc("/{module:.+}/@v/list", funcForList).Methods(http.MethodGet)
	router.HandleFunc("/{module:.+}/@v/{version}.info", funcForVersions).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(goAway)
	return router
}

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

func funcForList(w http.ResponseWriter, r *http.Request) {
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
