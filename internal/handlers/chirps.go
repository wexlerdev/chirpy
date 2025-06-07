package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/wexlerdev/chirpy/internal/database"
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
	fmt.Printf("input uuid for chirp: %v", chirpReq)
	if err != nil {
		respondWithError(w, 500, "error decoding req for chirp", err)
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

func mapDbChirpToChirp(dbChirp database.Chirp) *Chirp {
	return &Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt.Time,
		UpdatedAt: dbChirp.UpdatedAt.Time,
		Body:		dbChirp.Body,
		UserId:		dbChirp.UserID.UUID,
	}
}
