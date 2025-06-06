package handlers
import (
	"encoding/json"
	"net/http"
	"strings"
)


func (api *API) ValidateChirpHandler(w http.ResponseWriter, req * http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Decode params", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals {
		Body: cleanString(params.Body),
	})
}

func cleanString(str string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	const cleanWurd = "****"

	strSlice := strings.Split(str, " ")

	for _, badWurd := range profanity {
		for i, bodyWurd := range strSlice {
			if strings.ToLower(bodyWurd) == badWurd {
				strSlice[i] = cleanWurd
			}
		}
	}

	finalString := strings.Join(strSlice, " ")

	return finalString
}
