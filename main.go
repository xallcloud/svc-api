package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	gcp "github.com/xallcloud/gcp"
)

const (
	appName         = "svc-api"
	appVersion      = "0.0.2-demo"
	httpDefaultPort = "8080"
	topicPubNotify  = "notify"
	projectID       = "xallcloud"
)

// global resources for service
var dsClient *datastore.Client
var psClient *pubsub.Client
var tcPubNot *pubsub.Topic

func main() {
	// service initialization
	log.SetFlags(log.Lshortfile)

	log.Println("Starting", appName, "version", appVersion)

	port := os.Getenv("PORT")
	if port == "" {
		port = httpDefaultPort
		log.Printf("Service: %s. Defaulting to port %s", appName, port)
	}

	var err error
	ctx := context.Background()

	// DATASTORE initialization
	log.Println("Connect to Google 'datastore' on project: " + projectID)
	dsClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}

	// PUBSUB initialization
	log.Println("Connect to Google 'pub/sub' on project: " + projectID)
	psClient, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// topic to publish messages do type notify
	tcPubNot, err = gcp.CreateTopic(topicPubNotify, psClient)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}
	log.Printf("Using topic %v to post actions.\n", tcPubNot)

	// HTTP Server initialization
	// define all the routes for the HTTP server.
	//   The implementation is done on the "handler*.go" files
	router := mux.NewRouter()
	// version
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
	//generic options response
	router.HandleFunc("/api/action", optionsGeneric).Methods("OPTIONS")
	router.HandleFunc("/api/actions", optionsGeneric).Methods("OPTIONS")
	// Events
	router.HandleFunc("/api/events", getEvents).Methods("GET")
	router.HandleFunc("/api/events/callpoint/{cpID}", getEventsByCallpoint).Methods("GET")
	router.HandleFunc("/api/events/action/{acID}", getEventsByAction).Methods("GET")
	//generic options response
	router.HandleFunc("/api/events", optionsGeneric).Methods("OPTIONS")
	router.HandleFunc("/api/events/callpoint/{cpID}", optionsGeneric).Methods("OPTIONS")
	router.HandleFunc("/api/events/action/{acID}", optionsGeneric).Methods("OPTIONS")

	// Allow CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Start web server
	log.Printf("Service: %s. Listening on port %s", appName, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
