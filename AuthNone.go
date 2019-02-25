package middleware

import (
	"net/http"
)

// AuthSkip is a NOOP, it is useful when using AuthByScheme; skipping Known Schemes that will be validated later
func AuthSkip() func(h http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return next
	}
}
