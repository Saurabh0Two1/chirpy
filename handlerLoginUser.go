package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"
	"saurabh/chirpy.com/m/internal/database"
	"time"
)

func (cfg *apiConfig) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type UserLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userLoginRequest := UserLoginRequest{}
	err := decoder.Decode(&userLoginRequest)

	type errorResp struct {
		Error string `json:"error"`
	}

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

	userDetails, err := cfg.db.FindUserByEmail(r.Context(), userLoginRequest.Email)

	if err != nil {
		log.Printf("Failed to find user: %s", err)
		respBody := errorResp{
			Error: fmt.Sprintf("User with given email does not exist %s", err),
		}
		jsonData, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(jsonData)
		return
	}

	match, err := auth.CheckPasswordHash(userLoginRequest.Password, userDetails.HashedPassword)

	if match == false {
		log.Printf("wrong password: %s", err)
		respBody := errorResp{
			Error: fmt.Sprintf("Incorrect password. Please try again. %s", err),
		}
		jsonData, _ := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write(jsonData)
		return
	}

	token, err := auth.MakeJWT(userDetails.ID, cfg.secret)

	refTokenStr := auth.MakeRefreshToken()

	refToken := database.CreateRefreshTokenParams{
		Token:     refTokenStr,
		UserID:    userDetails.ID.String(),
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	}

	refresh_token, err := cfg.db.CreateRefreshToken(r.Context(), refToken)

	tz, err := time.LoadLocation("Asia/Kolkata")

	type successResp struct {
		ID           string `json:"id"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}

	userResp := successResp{
		Email:        userDetails.Email.String,
		CreatedAt:    userDetails.CreatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		UpdatedAt:    userDetails.UpdatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		ID:           userDetails.ID.String(),
		IsChirpyRed:  userDetails.IsChirpyRed,
		Token:        token,
		RefreshToken: refresh_token.Token,
	}

	jsonData, _ := json.Marshal(userResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)

}
