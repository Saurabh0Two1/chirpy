package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	signingKey := []byte(tokenSecret)

	expiresIn := time.Duration(3600*1000) * time.Millisecond

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Issuer:    string(TokenTypeAccess),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}
