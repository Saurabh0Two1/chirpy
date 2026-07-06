package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"time"

	"saurabh/chirpy.com/m/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request, authorUuid uuid.UUID) ([]database.Chirp, error) {
	if authorUuid != uuid.Nil {
		chirps, err := cfg.db.GetAllChirpsByUser(r.Context(), authorUuid)

		if err != nil {
			respondWithError(w, http.StatusBadRequest, "author ID is invalid", err)
			return nil, err
		}

		return chirps, nil
	}

	chirps, err := cfg.db.GetAllChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "author ID is invalid", err)
		return nil, err
	}

	return chirps, nil

}

func authorIDFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString == "" {
		return uuid.Nil, nil
	}
	authorID, err := uuid.Parse(authorIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (cfg *apiConfig) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	authorID, err := authorIDFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "author ID is invalid", err)
		return
	}

	chirps, err := cfg.GetChirps(w, r, authorID)

	sortDirection := r.URL.Query().Get("sort")

	slices.SortStableFunc(chirps, func(a, b database.Chirp) int {
		if sortDirection == "desc" {
			return b.CreatedAt.Compare(a.CreatedAt)
		}

		return a.CreatedAt.Compare(b.CreatedAt)
	})

	type successResp struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	var mappedChirps []successResp

	tz, _ := time.LoadLocation("Asia/Kolkata")

	for _, chirp := range chirps {
		mappedChirps = append(mappedChirps, successResp{
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
			ID:        chirp.ID.String(),
			UserID:    chirp.UserID.UUID.String(),
			UpdatedAt: chirp.UpdatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(mappedChirps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
