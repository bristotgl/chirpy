package main

import (
	"net/http"

	"github.com/bristotgl/chirpy/internal/auth"
)

func (cfg *apiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find refresh token", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), requestRefreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
