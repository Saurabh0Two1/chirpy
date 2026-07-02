package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"
	"saurabh/chirpy.com/m/internal/database"
	"time"
)

func (cfg *apiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type successResp struct {
		ID          string `json:"id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	type UserDataRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userDataRequest := UserDataRequest{}
	err := decoder.Decode(&userDataRequest)

	if err != nil {
		respBody := errorResp{
			Error: fmt.Sprintf("Failed to process request - %s", err),
		}

		jsonData, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(jsonData)
		return
	}

	if userDataRequest.Email == "" || userDataRequest.Password == "" {
		respBody := errorResp{
			Error: fmt.Sprintf("Email and Password are both required."),
		}
		jsonData, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(jsonData)
		return
	}

	hashedPwd, err := auth.HashPassword(userDataRequest.Password)

	if err != nil {
		respBody := errorResp{
			Error: fmt.Sprintf("Failed to preocess request."),
		}
		jsonData, _ := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(jsonData)
		return
	}

	user, _ := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: userDataRequest.Email, HashedPassword: hashedPwd})

	tz, err := time.LoadLocation("Asia/Kolkata")

	userResp := successResp{
		Email:       user.Email.String,
		CreatedAt:   user.CreatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		UpdatedAt:   user.UpdatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		ID:          user.ID.String(),
		IsChirpyRed: user.IsChirpyRed,
	}

	jsonData, _ := json.Marshal(userResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(jsonData)
}
