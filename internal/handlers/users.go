package handlers

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

func (api *API) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
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

	dbUser, err := api.cfg.DbQueries.CreateUser(req.Context(), emailRequest.Email)
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

func (api * API) ResetUsersHandler(w http.ResponseWriter, req *http.Request)  {
	//use devOnlyMiddleWare for auth, so assume here that it is authorized
	numDeleted, err := api.cfg.DbQueries.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, 500, "error deleting users", err)
	}

	type response struct {
		UsersDeleted	int64 `json:"users_deleted"`
	}

	res := response{
		UsersDeleted: numDeleted,
	}

	//reset fileServerHits
	api.cfg.FileserverHits.Store(0)

	respondWithJSON(w, 200, res)
}

func (api* API) DevOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if api.cfg.GetPlatform() != "dev" {
			http.Error(w, "Forbidden: this is only available during dev", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
