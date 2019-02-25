package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// AuthJWT validates a jwt and populates identity in context
//
// Example:
//   http.Handle("/", AuthJWT(fromValue, forKey, audience)(r))
func AuthJWT(
	fromValue func(r *http.Request) (string, error),
	forKey func(keyID string) []byte,
	audience string,
) func(h http.HandlerFunc) http.HandlerFunc {

	setTokenIdentity := func(ctx context.Context, token *jwt.Token) {
		// set token
		// set Identity from Claims
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		const KeyID string = "kid"

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			//Only accept expected valid signing methods - specifically, don't accept None
			log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("invalid method")
		}

		kid := ""
		if v, ok := token.Header[KeyID]; ok {
			kid = v.(string)
		}
		pem := forKey(kid)

		key, err := jwt.ParseRSAPublicKeyFromPEM(pem)
		if err != nil {
			return nil, fmt.Errorf("invalid key")
		}

		return key, nil
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			tokenString, err := fromValue(r)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				log.Printf("jwtValue: %+v\n", err)
				return
			}

			token, err := jwt.ParseWithClaims(
				tokenString,
				jwt.MapClaims{},
				keyFunc,
			)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				log.Printf("jwtParse: %+v\n", err)
				return
			}

			//TODO: Validate other required claims... issuer?

			mc := token.Claims.(jwt.MapClaims)
			if !mc.VerifyAudience(audience, true) {
				http.Error(w, "invalid audience", http.StatusUnauthorized)
				log.Printf("invalid audience, expected: %s, was: %v", audience, mc["aud"])
				return
			}

			setTokenIdentity(r.Context(), token)

			next.ServeHTTP(w, r)
		}
	}
}
