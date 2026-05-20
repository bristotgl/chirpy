package main

import (
	"fmt"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("./metrics.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error loading metrics page", err)
		return
	}

	respondWithText(w, http.StatusOK, fmt.Sprintf(string(data), cfg.fileserverHits.Load()), "text/html")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
