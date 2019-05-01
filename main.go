package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"

	_ "github.com/xallcloud/gcp"
)

const (
	appName    = "svc-api"
	appVersion = "0.0.1-alfa009"
	httpPort   = "8080"
	topicName  = "topicApi"
	projectID  = "xallcloud"
)

var dsClient *datastore.Client
var topic *pubsub.Topic

func main() {
	log.SetFlags(log.LstdFlags)
	log.Println("Starting", appName, "version", appVersion)

	port := os.Getenv("PORT")
	if port == "" {
		port = httpPort
		log.Printf("Service: %s. Defaulting to port %s", appName, port)
	}

	ctx := context.Background()
	// DATASTORE Initialization
	log.Println("Connect to Google datastore on project: " + projectID)
	var err error
	dsClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}

	router := mux.NewRouter()
	//define all the routes. The implementation is done on the "handlers.go" file
	router.HandleFunc("/api/version", getVersion).Methods("GET")
	//callpoints
	router.HandleFunc("/api/callpoints", getCallpoints).Methods("GET")
	router.HandleFunc("/api/callpoint/{id}", deleteCallpoint).Methods("DELETE")
	router.HandleFunc("/api/callpoint", postCallpoint).Methods("POST")
	//devices
	router.HandleFunc("/api/devices", getDevices).Methods("GET")
	router.HandleFunc("/api/device/{id}", deleteDevice).Methods("DELETE")
	router.HandleFunc("/api/device", postDevice).Methods("POST")

	// Start service
	log.Printf("Service: %s. Listening on port %s", appName, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
