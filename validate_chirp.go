package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request body", err)
		return
	}

	if len(strings.TrimSpace(params.Body)) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty", nil)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedText := cleanProfaneWords(params.Body)

	respondWithJSON(w, http.StatusOK, response{CleanedBody: cleanedText})
}

func cleanProfaneWords(text string) string {
	profaneWords := []string{"sharbert", "kerfuffle", "fornax"}

	textWords := strings.Split(text, " ")
	for i := range textWords {
		for _, profane := range profaneWords {
			if strings.ToLower(textWords[i]) == profane {
				textWords[i] = "****"
			}
		}
	}

	return strings.Join(textWords, " ")
}
