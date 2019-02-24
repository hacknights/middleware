package middleware

import (
	"context"
	"net/http"
)

type authenticatorFunc func(user, pass string) (map[string]interface{}, error)

// AuthBasic validates Basic credentials and populates identity in context
//
// Example:
//   http.Handle("/", AuthBasic(a.identity.ByCredentials)(r))
func AuthBasic(authenticate authenticatorFunc) func(h http.HandlerFunc) http.HandlerFunc {

	setIdentity := func(ctx context.Context, claims map[string]interface{}) {

	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			//TODO: is this required?
			w.Header().Set("WWW-Authenticate", `Basic realm="Fail"`)

			user, pass, ok := r.BasicAuth()
			if !ok {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

			claims, err := authenticate(user, pass)
			if err != nil {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

			setIdentity(r.Context(), claims)

			next.ServeHTTP(w, r)
		}
	}
}
