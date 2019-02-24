package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recuperate catches panics in handlers, logs the stack trace and serves an HTTP 500 error.
//
// Example:
//   http.Handle("/", Recuperate(r))
func Recuperate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "HTTP 500: internal server error (we've logged it!)", http.StatusInternalServerError)
				log.Printf("Handler panicked: %s\n%s", err, debug.Stack())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
