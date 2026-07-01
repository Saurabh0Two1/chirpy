package main

import (
	"database/sql"
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"
	"saurabh/chirpy.com/m/internal/database"
	"time"
)

func (cfg *apiConfig) RevokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil || refreshToken == "" {
		respondWithError(w, 401, "refresh token is missing", err)
	}

	revokeArg := database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: time.Now(),
		Token:     refreshToken,
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), revokeArg)

	if err != nil {
		respondWithError(w, 401, "refresh token is missing", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)

}
