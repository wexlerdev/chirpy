package handlers

import (
	"time"
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"errors"
	"github.com/wexlerdev/chirpy/internal/database"
	"github.com/wexlerdev/chirpy/internal/auth"
	"fmt"
)


type User struct {
	ID        uuid.NullUUID `json:"id"`
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

	user.ID = uuid.NullUUID {
		UUID: dbUser.ID,
		Valid: true,
	}
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
		return
	}
	//check password
	err = auth.CheckPasswordHash(dbUser.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "password not correct", err)
		return
	}

	user := mapDbUserToUser(dbUser)

	defaultTokenExpiration := 3600 * time.Second //1 hr
	//make jwt (token)
	token, err := auth.MakeJWT(user.ID.UUID, api.cfg.JwtSecret, defaultTokenExpiration) 
	if err != nil {
		respondWithError(w, 500, "issue making token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "issue making refresh token", err)
		return
	}


	refreshTokenParams := database.CreateRefreshTokenParams {
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * (24 * 60)),
	}

	_, err = api.cfg.DbQueries.CreateRefreshToken(req.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(w, 500, "error creating refresh token", err)
		return
	}

	responseObj := struct {
		User
		Token	string `json:"token"`
		RefreshToken	string	`json:"refresh_token"`
	} {
		User: *user,
		Token: token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, 200, responseObj)
}

func (api *API) HandleRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "error getiing token from header", err)
		return
	}
	//look up refresh token in db	
	refreshTokenStruct, err := api.cfg.DbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "cannot find refresh token", err)
		return
	}

	if time.Now().After(refreshTokenStruct.ExpiresAt) {
		respondWithError(w, 401, "token has expired", errors.New("token has expired"))
		return
	}
	fmt.Printf("refreshTokenRevokedVal: %v, Valid: %v \n", refreshTokenStruct.RevokedAt, refreshTokenStruct.RevokedAt.Valid)

	//if the refreshToken was revoked and that it is after the revoked time
	if refreshTokenStruct.RevokedAt.Valid && time.Now().After(refreshTokenStruct.RevokedAt.Time) {
		respondWithError(w, 401, "token has expired", errors.New("token has expired"))
		return
	}

	//get user from db using refresh token
	user, err := api.cfg.DbQueries.GetUserFromRefreshToken(req.Context(), refreshTokenStruct.Token)
	if err != nil {
		respondWithError(w, 401, "could not find user from refresh token", err)
		return
	}
	//make jwt using uuid from user
	token, err := auth.MakeJWT(user.ID, api.cfg.JwtSecret, time.Second * 3600)
	if err != nil {
		respondWithError(w, 401, "error creating jwt", err)
		return
	}
	//respond with dat token!
	responseStruct := struct {
		Token	string `json:"token"`
	} {
		Token: token,
	}
	respondWithJSON(w, 200, responseStruct)
}

func (api *API) HandleRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "error getiing token from header", err)
		return
	}
	refreshTok, err := api.cfg.DbQueries.GetRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, 401, "error getting refreshTok from db", err)
		return
	}

	//now lets revoke the tok in the database homie
	err = api.cfg.DbQueries.RevokeRefreshToken(req.Context(), refreshTok.Token)
	if err != nil {
		respondWithError(w, 401, "error revoking the token in db", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
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

