package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

type errorResp struct {
	Error string `json:"error"`
}

type chirpValidStatus struct {
	Valid bool `json:"valid"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	chirp := chirp{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request body", err)
		return
	}

	if len(strings.TrimSpace(chirp.Body)) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty", nil)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	resp := chirpValidStatus{Valid: true}
	encodedResp, err := json.Marshal(resp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error marshaling response into JSON", err)
		return
	}

	writeResponse(w, http.StatusOK, encodedResp, "application/json")
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	log.Printf("%s: %s", msg, err)

	errorResp := errorResp{Error: msg}
	encodedErrorResp, err := json.Marshal(errorResp)
	if err != nil {
		payload := []byte("{error: Unexpected failure happened.}")
		writeResponse(w, http.StatusInternalServerError, payload, "application/json")
		return
	}

	writeResponse(w, code, encodedErrorResp, "application/json")
}

func writeResponse(w http.ResponseWriter, code int, payload any, contentType string) {
	w.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(code)

	if payload != nil {
		w.Write(payload.([]byte))
	}
}
