package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/bristotgl/chirpy/internal/auth"
	"github.com/bristotgl/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	requestRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), requestRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token does not exist", err)
		return
	}

	if err := ValidateRefreshToken(refreshToken); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is invalid", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user for refresh token", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}

func ValidateRefreshToken(refreshToken database.RefreshToken) error {
	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		return errors.New("Refresh token is expired")
	}

	if refreshToken.RevokedAt.Valid {
		return errors.New("Refresh token is revoked")
	}

	return nil
}
