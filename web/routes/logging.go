package routes

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"service": "web"})

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
	return
}

func loggingWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)
		h.ServeHTTP(wrapped, r)
		log.WithFields(logrus.Fields{
			"status":   wrapped.status,
			"method":   r.Method,
			"path":     r.URL.EscapedPath(),
			"duration": time.Since(start),
		}).Info("Served request")
	})
}
