package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api", apiHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, `{"name": "xallcloud", "version": "0.1.0"}`)
}
