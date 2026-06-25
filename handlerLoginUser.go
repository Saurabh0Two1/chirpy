package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saurabh/chirpy.com/m/internal/auth"
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
	fmt.Printf("a == %v \n", userLoginRequest)

	userDetails, err := cfg.db.FindUserByEmail(r.Context(), userLoginRequest.Email)

	fmt.Printf("b == %v \n", userDetails)

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

	fmt.Printf("p== %v %v \n", userLoginRequest.Password, userDetails.HashedPassword)

	match, err := auth.CheckPasswordHash(userLoginRequest.Password, userDetails.HashedPassword)

	fmt.Printf("q== %v %v \n", match, err)

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

	tz, err := time.LoadLocation("Asia/Kolkata")

	type successResp struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
	}

	userResp := successResp{
		Email:     userDetails.Email.String,
		CreatedAt: userDetails.CreatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		UpdatedAt: userDetails.UpdatedAt.In(tz).Format("2006-01-02T15:04:05 +05:30:00"),
		ID:        userDetails.ID.String(),
	}

	jsonData, _ := json.Marshal(userResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)

}
