package handlers

import (
	"time"
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"github.com/wexlerdev/chirpy/internal/database"
	"github.com/wexlerdev/chirpy/internal/auth"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (api *API) CreateUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password	string `json:"password"`
	}
	var params parameters

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "error decoding req email", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "error hashing password", err)
		return
	}	

	dbUser, err := api.cfg.DbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})

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

func (api * API) HandleLogin(w http.ResponseWriter, req * http.Request) {
	type parameters struct {
		Password string
		Email string
	}
	var params parameters

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "error parsing login params", err)
		return
	}
	//find user by email
	dbUser, err := api.cfg.DbQueries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "error finding user by email", err)
	}
	//check password
	err = auth.CheckPasswordHash(dbUser.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "password not correct", err)
	}

	user := mapDbUserToUser(dbUser)
	respondWithJSON(w, 200, user)
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
