package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bristotgl/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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

	cleanedBody := cleanProfaneWords(params.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
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
