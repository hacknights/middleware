package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// AuthJWT validates a jwt and populates identity in context
//
// Example:
//   http.Handle("/", jwtAuth(jwtValueFunc, jwt.Keyfunc)(r))
func AuthJWT(key interface{}) func(h http.HandlerFunc) http.HandlerFunc {

	setTokenIdentity := func(ctx context.Context, token *jwt.Token) {
		// set token
		// set Identity from Claims
	}

	jwtValue := func(r *http.Request) (string, error) {
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			return "", fmt.Errorf("missing authorization header")
		}
		return s[1], nil
	}

	keyFunc := func(key interface{}) func(token *jwt.Token) (interface{}, error) {
		return func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				//Only accept expected valid signing methods - specifically, don't accept None
				log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
				return nil, fmt.Errorf("invalid method")
			}
			//TODO: use token.Header["kid"]
			return key, nil
		}
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			tokenString, err := jwtValue(r)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				log.Printf("jwtValue: %+v\n", err)
				return
			}

			token, err := jwt.Parse(
				tokenString,
				keyFunc(key),
			)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				log.Printf("jwtParse: %+v\n", err)
				return
			}

			//TODO: Validate required claims: Issuer (iss), Audience (aud)

			setTokenIdentity(r.Context(), token)

			next.ServeHTTP(w, r)
		}
	}
}
