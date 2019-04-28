package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	appName    = "svc-api"
	appVersion = "0.0.1-alfa001"
	httpPort   = "8080"
	topicName  = "topicNotifification"
	projectID  = "xallcloud"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = httpPort
		log.Printf("Service: %s. Defaulting to port %s", appName, port)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/version", getVersionHanlder).Methods("GET")
	router.HandleFunc("/api/callpoints", getCallpointsHanlder).Methods("GET")

	log.Printf("Service: %s. Listening on port %s", appName, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func getVersionHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("[/version:GET] Requested api version.")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, fmt.Sprintf(`{"service": "%s", "version": "%s"}`, appName, appVersion))
}

func getCallpointsHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoints:GET] Requested get all callpoints")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"callpoints": {} }`)
}
