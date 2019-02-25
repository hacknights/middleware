package middleware

import (
	"net/http"
	"strings"
)

// AuthByScheme delegates to the appropriate middleware based on the HTTP Authorization Scheme
//
// Example:
//   http.Handle("/", AuthByScheme(AuthBasic(...), AuthJwt(...))(r))
func AuthByScheme(schemeHandlers map[string]func(h http.HandlerFunc) http.HandlerFunc) func(h http.HandlerFunc) http.HandlerFunc {

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			a := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			if len(a) < 2 {
				http.Error(w, "invalid authorization scheme", http.StatusUnauthorized)
				return
			}
			value := a[0]

			for scheme, h := range schemeHandlers {
				if strings.EqualFold(value, scheme) {
					h(next)(w, r)
					return
				}
			}

			http.Error(w, "invalid authorization scheme", http.StatusUnauthorized)
		}
	}
}
