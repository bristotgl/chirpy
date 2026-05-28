package main

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/bristotgl/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	var databaseChirps []database.Chirp
	var err error
	var errorStatus int

	userID := r.URL.Query().Get("author_id")
	if userID != "" {
		databaseChirps, errorStatus, err = cfg.getAllChirpsOfUser(r, userID)
	} else {
		databaseChirps, errorStatus, err = cfg.getAllChirps(r)
	}

	if err != nil {
		respondWithError(w, errorStatus, "Error getting chirps", err)
		return
	}

	sortChirps(databaseChirps, r.URL.Query().Get("sort"))

	response := []Chirp{}
	for _, chirp := range databaseChirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func sortChirps(chirps []database.Chirp, order string) {
	switch order {
	case "asc":
		slices.SortFunc(chirps, sortChirpsAsc)
	case "desc":
		slices.SortFunc(chirps, sortChirpsDesc)
	}
}

func sortChirpsAsc(a, b database.Chirp) int {
	if a.CreatedAt.Before(b.CreatedAt) {
		return -1
	}
	if a.CreatedAt.After(b.CreatedAt) {
		return 1
	}
	return 0
}

func sortChirpsDesc(a, b database.Chirp) int {
	if a.CreatedAt.After(b.CreatedAt) {
		return -1
	}

	if a.CreatedAt.Before(b.CreatedAt) {
		return 1
	}
	return 0
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) getAllChirpsOfUser(r *http.Request, userID string) ([]database.Chirp, int, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return []database.Chirp{}, http.StatusBadRequest, fmt.Errorf("Error parsing user ID: %w", err)
	}

	chirps, err := cfg.db.GetAllChirpsOfUser(r.Context(), parsedUserID)
	if err != nil {
		return []database.Chirp{}, http.StatusNotFound, fmt.Errorf("Chirp not found for user '%s': %w", userID, err)
	}

	return chirps, 0, nil
}

func (cfg *apiConfig) getAllChirps(r *http.Request) ([]database.Chirp, int, error) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		return []database.Chirp{}, http.StatusInternalServerError, fmt.Errorf("Error getting all chirps: %w", err)
	}

	return chirps, 0, nil
}
