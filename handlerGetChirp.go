package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	rawID := r.PathValue("chirpID")

	// 2. Parse the string to validate it as a UUID
	chirpID, err := uuid.Parse(rawID)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)

	if err != nil {
		http.Error(w, "Chirp not found", http.StatusNotFound)
		return
	}

	type successResp struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	tz, err := time.LoadLocation("Asia/Kolkata")

	respBody := successResp{
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		ID:        chirp.ID.String(),
		UserID:    chirp.UserID.UUID.String(),
		UpdatedAt: chirp.UpdatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
	}

	jsonData, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)

}
