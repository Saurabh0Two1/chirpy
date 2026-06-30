package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	headerParts := strings.Split(authHeader, " ")
	fmt.Printf("headers %v \n", headerParts)
	token := headerParts[1]
	return token, nil
}
