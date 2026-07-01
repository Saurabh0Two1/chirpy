package main

import (
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	// authenticate user
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Please login to proceed", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)

	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid user credentials", err)
		return

	}

	rawID := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(rawID)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp Id is not correct", err)
		return

	}

	chirpData, err := cfg.db.GetChirp(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return

	}

	chirpUserId, err := uuid.Parse(chirpData.UserID.UUID.String())

	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found for the chirp", err)
		return

	}

	if chirpUserId != userId {
		respondWithError(w, http.StatusForbidden, "This is not your chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)

	if err != nil {
		respondWithError(w, http.StatusForbidden, "User not found for the chirp", err)
		return

	}

	respondWithJSON(w, http.StatusNoContent, nil)
	return

	//make delete db request

}
