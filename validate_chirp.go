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
		Valid bool `json:"valid"`
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

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error marshaling response into JSON", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Valid: true})
}
