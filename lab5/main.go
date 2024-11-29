package main

import (
	"log"
	"net/http"
)

const (
	ADDRESS  = "localhost:8080"
	BASE_DIR = "./public"
)

func main() {
	fileServer := http.FileServer(http.Dir(BASE_DIR))
	http.Handle("/", fileServer)
	log.Printf("Serving files from %s on http://%s", BASE_DIR, ADDRESS)
	err := http.ListenAndServe(ADDRESS, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
