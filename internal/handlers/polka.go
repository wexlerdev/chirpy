package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/wexlerdev/chirpy/internal/auth"
)

func (api *API) HandleUserUpgrade(w http.ResponseWriter, req * http.Request) {
	type dataStruct struct {
		UserId		uuid.UUID	`json:"user_id"`
	}

	type input struct {
		Event		string		`json:"event"`
		Data		dataStruct	`json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, 401, "err getting apikey", err)
		return
	}

	if apiKey != api.cfg.PolkaKey {
		respondWithError(w, 401, "apiKey unauthorized", err)
		return
	}

	var inputStruct input

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err = decoder.Decode(&inputStruct)
	if err != nil {
		respondWithError(w, 400, "err decoding body", err)
		return
	}
	//if input.Event is not user.upgraded, repond with 204
	if inputStruct.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	_, err = api.cfg.DbQueries.UpgradeUser(req.Context(), inputStruct.Data.UserId)
	if err != nil {
		respondWithError(w,404, "did not find user", err)
		return
	}
	w.WriteHeader(204)
}
