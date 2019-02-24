package middleware

import "net/http"

// Use wraps a handler func with additional behaviors
//
// Example:
//   http.Handle("/", Use(handler, TraceIDs, Logging, Recuperate)
func Use(middleware ...func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		for _, m := range middleware {
			h = m(h)
		}
		return h
	}
}
