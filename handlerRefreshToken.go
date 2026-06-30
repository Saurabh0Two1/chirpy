package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"
	"time"
)

func (cfg *apiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil || refreshToken == "" {
		returnError(w, 401, "refresh token is missing", err)
	}

	refreshTokenData, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)

	fmt.Printf("\n expire: %v \n revoke: %v \n time now: %v \n  err: %v \n revoke valid %v \n revoke now: %v \n revoke before: %v \n",
		refreshTokenData.ExpiresAt.Before(time.Now()),
		refreshTokenData.RevokedAt.Time.UTC(),
		time.Now().UTC(),
		err,
		refreshTokenData.RevokedAt.Valid,
		refreshTokenData.RevokedAt.Time.UTC().Equal(time.Now().UTC()),
		refreshTokenData.RevokedAt.Time.UTC().Before(time.Now().UTC()))

	if err != nil {
		returnError(w, 401, "failed to find user", err)
		return
	}

	if refreshTokenData.ExpiresAt.Before(time.Now()) {
		returnError(w, 401, "token expired", err)
		return
	}

	if refreshTokenData.RevokedAt.Valid && refreshTokenData.RevokedAt.Time.Before(time.Now()) {
		returnError(w, 401, "token revoked", err)
		return
	}

	accessToken, err := auth.MakeJWT(refreshTokenData.UserID.UUID, cfg.secret)

	if err != nil {
		returnError(w, 500, "failed to create access token", err)
		return
	}

	type successResp struct {
		Token string `json:"token"`
	}

	resp := successResp{
		Token: accessToken,
	}

	jsonData, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}
