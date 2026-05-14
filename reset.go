package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	respondWithText(w, http.StatusOK, "Hits reset to 0", "text/plain")
}
