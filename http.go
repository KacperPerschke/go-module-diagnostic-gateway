package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	const healthContent = `
    <!DOCTYPE html>
    <html>
        <body>
            <p>Pocałujta w … wójta!</p>
        </body>
    </html>
    `
	fmt.Fprintf(w, healthContent)
}

func handlerForPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `pong`)
}

func produceRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	globalLogger.Info("Starts preparation of http routing.")
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
	globalLogger.Info("Successfully finished preparation of http routing.")
	return router
}

/*
 * type (
 * 	// struct for holding response details
 * 	responseData struct {
 * 		status int
 * 		size   int
 * 	}
 *
 * 	// our http.ResponseWriter implementation
 * 	loggingResponseWriter struct {
 * 		http.ResponseWriter // compose original http.ResponseWriter
 * 		responseData        *responseData
 * 	}
 * )
 *
 * func (r *loggingResponseWriter) Write(b []byte) (int, error) {
 * 	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
 * 	r.responseData.size += size            // capture size
 * 	return size, err
 * }
 *
 * func (r *loggingResponseWriter) WriteHeader(statusCode int) {
 * 	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
 * 	r.responseData.status = statusCode       // capture status code
 * }
 *
 * func WithLogging(h http.Handler) http.Handler {
 * 	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
 * 		start := time.Now()
 *
 * 		responseData := &responseData{
 * 			status: 0,
 * 			size:   0,
 * 		}
 * 		lrw := loggingResponseWriter{
 * 			ResponseWriter: rw, // compose original http.ResponseWriter
 * 			responseData:   responseData,
 * 		}
 * 		h.ServeHTTP(&lrw, req) // inject our implementation of http.ResponseWriter
 *
 * 		duration := time.Since(start)
 *
 * 		logrus.WithFields(logrus.Fields{
 * 			"uri":      req.RequestURI,
 * 			"method":   req.Method,
 * 			"status":   responseData.status,
 * 			"duration": duration,
 * 			"size":     responseData.size,
 * 		}).Info("request completed")
 * 	}
 * 	return http.HandlerFunc(loggingFn)
 * }
 */
