package main

import "net/http"

func main() {
	mux := new(http.ServeMux)
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
