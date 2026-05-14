package main

import (
	"fmt"
	"net/http"
)

func respondWithText(w http.ResponseWriter, code int, payload string, contentType string) {
	w.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(code)
	w.Write([]byte(payload))
}
