package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

// traceIDs writes the RequestID, and possibly CorrelationID, in the Request header
//
// Example:
//   http.Handle("/", traceIDs(r))
func TraceIDs(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := new(http.Request)
		*r2 = *r
		rid := uuid.New().String()
		r2.Header.Set("X-Request-ID", rid)
		//TODO: if not already set, copy RequestID to CorrelationID...?
		next(w, r2)
	})
}
