package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	writeResponse(w, http.StatusOK, []byte("Hits reset to 0"), "text/plain")
}
