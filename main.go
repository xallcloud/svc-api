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

	gcp "github.com/xallcloud/gcp"
)

const (
	appName    = "svc-api"
	appVersion = "0.0.1.alfa.22-events-get"
	httpPort   = "8080"
	topicName  = "notify"
	projectID  = "xallcloud"
)

var dsClient *datastore.Client
var psClient *pubsub.Client
var topic *pubsub.Topic

func main() {
	log.SetFlags(log.LstdFlags)
	log.Println("Starting", appName, "version", appVersion)

	port := os.Getenv("PORT")
	if port == "" {
		port = httpPort
		log.Printf("Service: %s. Defaulting to port %s", appName, port)
	}

	var err error
	ctx := context.Background()
	// DATASTORE Initialization
	log.Println("Connect to Google 'datastore' on project: " + projectID)

	dsClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}

	// PUBSUB Initialization
	log.Println("Connect to Google 'pub/sub' on project: " + projectID)

	psClient, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	topic, err = gcp.CreateTopic(topicName, psClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using topic %v to post actions.\n", topic)

	// HTTP SERVER Initialization
	router := mux.NewRouter()
	// define all the routes for the HTTP server.
	// The implementation is done on the "handler*.go" files
	// Common to all services
	router.HandleFunc("/api/version", getVersion).Methods("GET")
	// Callpoints
	router.HandleFunc("/api/callpoints", getCallpoints).Methods("GET")
	router.HandleFunc("/api/callpoint/{id}", deleteCallpoint).Methods("DELETE")
	router.HandleFunc("/api/callpoint", postCallpoint).Methods("POST")
	// Devices
	router.HandleFunc("/api/devices", getDevices).Methods("GET")
	router.HandleFunc("/api/device/{id}", deleteDevice).Methods("DELETE")
	router.HandleFunc("/api/device", postDevice).Methods("POST")
	// Assignments
	router.HandleFunc("/api/assignment", postAssignment).Methods("POST")
	router.HandleFunc("/api/assignments/callpoint/{cpID}", getAssignmentsByCallpoint).Methods("GET")
	// Actions
	router.HandleFunc("/api/actions", getActions).Methods("GET")
	router.HandleFunc("/api/action", postAction).Methods("POST")
	//Events
	router.HandleFunc("/api/events", getEvents).Methods("GET")
	router.HandleFunc("/api/events/callpoint/{cpID}", getEventsByCallpoint).Methods("GET")

	// Start service
	log.Printf("Service: %s. Listening on port %s", appName, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
