package main

import (
	"encoding/json"
	"net/http"

	"saurabh/chirpy.com/m/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) UpgradeUserHandler(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "failed to get the API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "auth header key does not match the actual key", err)
		return
	}

	type UpgradeRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	upgradeRequest := UpgradeRequest{}
	err = decoder.Decode(&upgradeRequest)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to parse upgrade request", err)
		return
	}

	if upgradeRequest.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	if upgradeRequest.Event == "user.upgraded" {
		userId, err := uuid.Parse(upgradeRequest.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to parse the user ID", err)
			return
		}

		userData, err := cfg.db.FindUserByID(r.Context(), userId)

		if err != nil {
			respondWithError(w, http.StatusNotFound, "failed to upgrade user", err)
			return
		}

		_, err = cfg.db.UpgradeToRed(r.Context(), userData.ID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to upgrade user", err)
			return
		}

		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

}
