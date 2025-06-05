package main

import (
	"time"
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/wexlerdev/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type emailReq struct {
		Email string `json:"email"`
	}
	var emailRequest emailReq


	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&emailRequest)
	if err != nil {
		respondWithError(w, 500, "error decoding req email", err)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), emailRequest.Email)
	if err != nil {
		respondWithError(w, 500, "error creating user", err)
		return
	}

	user := mapDbUserToUser(dbUser)
	respondWithJSON(w, 201, user)
}

func mapDbUserToUser(dbUser database.User) * User {
	var user User

	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email

	return & user
}
