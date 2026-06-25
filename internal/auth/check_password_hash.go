package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Failed to match passwords")
	}

	return match, nil
}
