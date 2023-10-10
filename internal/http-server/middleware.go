package httpserver

import (
	"fmt"
	"net/http"
	"time"
)

func (serv *_HTTPServer) LogMiddleware(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			status:         0,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		serv.Logger.Info(fmt.Sprintf("uri: %s; method: %s; status: %d; duration: %s", r.RequestURI, r.Method, lw.status, duration.String()))
	}

	return http.HandlerFunc(logFn)
}
