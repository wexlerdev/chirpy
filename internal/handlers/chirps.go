package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/wexlerdev/chirpy/internal/database"
	"github.com/wexlerdev/chirpy/internal/auth"
	"time"
	"fmt"
	"errors"
	
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

func (api * API) HandleDeleteChirp(w http.ResponseWriter, req * http.Request) {
	//get the access token in header
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 401, "did not get token from header", err)
		return
	}
	//get the chirpId from the path param
	idString := req.PathValue("id")
	idParam, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, 401, "did not parse idparam to uuid", err)
		return
	}
	//get the chirp from the db
	chirp, err := api.cfg.DbQueries.GetChirp(req.Context(), idParam)
	if err != nil {
		respondWithError(w, 404, "did not find chirp", err)
		return
	}

	//get the id in the jwt
	idJwt, err := auth.ValidateJWT(token, api.cfg.JwtSecret)
	if err != nil {
		respondWithError(w, 401, "did not validate jwt", err)
		return
	}

	//check if userId in chirp matching user id of access token
	if chirp.UserID.UUID != idJwt {
		respondWithError(w, 403, "user cannot delete this chirp", errors.New("user unauthorized"))
		return
	}
	//user is authorized :)
	_, err = api.cfg.DbQueries.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, 500, "error deleting chirp", err)
	}
	w.WriteHeader(204)
}

func (api * API) HandleGetAllChirps(w http.ResponseWriter, req * http.Request) {
	dbChirps, err := api.GetChirps(req)
	if err != nil {
		respondWithError(w, 400, "error getting chirps from db", err)	
	}

	chirps := make([]Chirp,0, len(*dbChirps))

	for _, dbChirp := range *dbChirps {
		chirps = append(chirps, *mapDbChirpToChirp(dbChirp))
	}

	respondWithJSON(w, 200, chirps)
}

func (api * API) GetChirps(req * http.Request) (*[]database.Chirp, error) {
	authorId := req.URL.Query().Get("author_id")
	sortString := req.URL.Query().Get("sort")

	if authorId != "" {
		authorUUID, err := uuid.Parse(authorId)
		if err != nil {
			return nil, err
		}
		authorNUUID := uuid.NullUUID{UUID: authorUUID, Valid: true}
		if sortString == "desc" {
			dbChirps, err := api.cfg.DbQueries.GetChirpsByAuthorIdDesc(req.Context(), authorNUUID)
			return &dbChirps, err
		}
		
		dbChirps, err := api.cfg.DbQueries.GetChirpsByAuthorIdAsc(req.Context(), authorNUUID)
		return &dbChirps, err
	} 

	if sortString == "desc" {
		dbChirps, err := api.cfg.DbQueries.GetAllChirpsDesc(req.Context())
		return &dbChirps, err
	}

	dbChirps, err := api.cfg.DbQueries.GetAllChirpsAsc(req.Context())
	return &dbChirps, err
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


