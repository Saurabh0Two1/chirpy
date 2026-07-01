package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) < 2 || authHeaderParts[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}
	token := authHeaderParts[1]
	return token, nil
}
