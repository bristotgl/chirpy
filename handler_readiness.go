package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithText(w, http.StatusOK, http.StatusText(http.StatusOK), "text/plain")
}
