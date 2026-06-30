package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *apiConfig) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirps, _ := cfg.db.GetAllChirps(r.Context())

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
