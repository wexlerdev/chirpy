package handlers
import (
	"errors"
	"strings"
)


func validateChirp(bod string) (string, error) {


	const maxChirpLength = 140
	if len(bod) > maxChirpLength {
		return bod, errors.New("chirp is too long")
	}

	bod = cleanString(bod)

	return bod, nil
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
