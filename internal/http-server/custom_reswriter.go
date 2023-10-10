package httpserver

import "net/http"

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	size, err := r.ResponseWriter.Write(b)
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}
