package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) < 2 || authHeaderParts[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	key := strings.Trim(authHeaderParts[1], " ")
	return key, nil
}
