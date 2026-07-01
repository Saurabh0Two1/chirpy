package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"
	"saurabh/chirpy.com/m/internal/database"
	"strings"
	"time"
)

func (cfg *apiConfig) CreateChirpHandler(w http.ResponseWriter, r *http.Request) {

	type errorResp struct {
		Error string `json:"error"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respBody := errorResp{
			Error: fmt.Sprintf("Error decoding token: %s", err),
		}
		jsonData, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(jsonData)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secret)

	if err != nil {
		respBody := errorResp{
			Error: fmt.Sprintf("Invalid token. Please re-login %v", err),
		}
		jsonData, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(jsonData)
		return
	}

	type Chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err = decoder.Decode(&chirp)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respBody := errorResp{
			Error: fmt.Sprintf("Error decoding parameters: %s", err),
		}
		jsonData, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(jsonData)
		return
	}

	if len(chirp.Body) > 140 {
		respBody := errorResp{
			Error: "Chirp is too long",
		}
		w.Header().Set("Content-Type", "application/json")
		jsonData, _ := json.Marshal(respBody)

		w.WriteHeader(400)
		w.Write(jsonData)
		return
	}

	words := strings.Split(chirp.Body, " ")

	for i := 0; i < len(words); i++ {
		switch strings.ToLower(words[i]) {
		case "kerfuffle":
			words[i] = "****"

		case "sharbert":
			words[i] = "****"

		case "fornax":
			words[i] = "****"

		}
	}

	body := strings.Join(words, " ")

	chirp1 := database.CreateChirpParams{
		Body:   body,
		UserID: userId.String(),
	}

	savedChirp, _ := cfg.db.CreateChirp(r.Context(), chirp1)

	type successResp struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	tz, err := time.LoadLocation("Asia/Kolkata")

	respBody := successResp{
		Body:      savedChirp.Body,
		CreatedAt: savedChirp.CreatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		ID:        savedChirp.ID.String(),
		UserID:    savedChirp.UserID.UUID.String(),
		UpdatedAt: savedChirp.UpdatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
	}

	jsonData, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonData)
}
