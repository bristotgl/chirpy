package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, []byte(http.StatusText(http.StatusOK)), "text/plain")
}
