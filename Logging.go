package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging writes the RequestURI and duration of handlers
//
// Example:
//   http.Handle("/", Logging(r))
func Logging(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := newLoggingResponseWriter(w)
		start := time.Now()
		defer func() {
			elapsed := time.Since(start).String()
			rid := r.Header["X-Request-Id"]
			log.Printf("%s %3d %7d %12s %s %s", rid, l.statusCode, l.length, elapsed, r.UserAgent(), r.RequestURI)
		}()
		next.ServeHTTP(l, r)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	length     int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		length:         0}
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := l.ResponseWriter.Write(b)
	l.length += n
	return n, err
}
