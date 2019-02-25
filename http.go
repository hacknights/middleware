package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func FromHeader(r *http.Request) (string, error) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return "", fmt.Errorf("missing authorization header")
	}
	return s[1], nil
}

func FromCookie(cookieName string) func(r *http.Request) (string, error) {
	return func(r *http.Request) (string, error) {
		return "", nil
	}
}
