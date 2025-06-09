package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/wexlerdev/chirpy/internal/database"
	"github.com/wexlerdev/chirpy/internal/auth"
	"time"
	"fmt"
	
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body	  string  `json:"body"`
	UserId  	uuid.UUID `json:"user_id"`
}


func (api *API) HandleCreateChirp(w http.ResponseWriter, req *http.Request) {
	type chirpRequest struct {
		Body	string		`json:"body"`
		UserId	uuid.NullUUID	`json:"user_id"`
	}
	var chirpReq chirpRequest

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&chirpReq)

	if err != nil {
		respondWithError(w, 500, "error decoding req for chirp", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	fmt.Printf("token: %v", token)
	if err != nil {
		respondWithError(w, 401, "error getting token", err)
		return
	}

	//authenticate user
	userId, err := auth.ValidateJWT(token, api.cfg.JwtSecret)
	chirpReq.UserId = uuid.NullUUID{
		UUID: userId,
		Valid: true,
	}

	if err != nil {
		respondWithError(w,401, "error validating token", err)
		return
	}

	//validate chirp
	chirpReq.Body, err = validateChirp(chirpReq.Body)
	if err != nil {
		respondWithError(w, 400, "chirp invalid", err)
		return
	}
	
	chirpParams := database.CreateChirpParams{
		Body: chirpReq.Body,
		UserID: chirpReq.UserId,
	}

	dbChirp, err := api.cfg.DbQueries.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		respondWithError(w, 500, "error creating chirp", err)
		return
	}

	chirp := mapDbChirpToChirp(dbChirp)
	respondWithJSON(w, 201, chirp)
}

func (api * API) HandleGetAllChirps(w http.ResponseWriter, req * http.Request) {
	dbChirps, err := api.cfg.DbQueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "error getting chirps from db", err)	
	}

	chirps := make([]Chirp, 0, len(dbChirps))

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, *mapDbChirpToChirp(dbChirp))
	}
	respondWithJSON(w, 200, chirps)
}

func (api * API) HandleGetChirp(w http.ResponseWriter, req * http.Request) {
	chirpId := req.PathValue("id")
	chirpUuid, err := uuid.Parse(chirpId) 
	if err != nil {
		respondWithError(w, 400, "error parsing chirp id", err)
		return
	}

	dbChirp, err := api.cfg.DbQueries.GetChirp(req.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, 404, "could not find chirp", err)
		return
	}
	respondWithJSON(w, 200, *mapDbChirpToChirp(dbChirp))
}

func mapDbChirpToChirp(dbChirp database.Chirp) *Chirp {
	return &Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt.Time,
		UpdatedAt: dbChirp.UpdatedAt.Time,
		Body:		dbChirp.Body,
		UserId:		dbChirp.UserID.UUID,
	}
}


