package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"

	pbn "github.com/xallcloud/api/proto"
	_ "github.com/xallcloud/gcp"
)

const (
	appName    = "svc-api"
	appVersion = "0.0.1-alfa001"
	httpPort   = "8080"
	topicName  = "topicApi"
	projectID  = "xallcloud"
)

var dsClient *datastore.Client
var psClient *pubsub.Client
var topic *pubsub.Topic

func main() {
	log.SetFlags(log.LstdFlags)
	log.Println("Starting", appName)

	port := os.Getenv("PORT")
	if port == "" {
		port = httpPort
		log.Printf("Service: %s. Defaulting to port %s", appName, port)
	}

	ctx := context.Background()
	// DATASTORE
	log.Println("Connect to Google datastore on project: " + projectID)
	var err error
	dsClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}
	// PUB/SUB
	log.Println("Connect to Google Pub/Sub on project: " + projectID)
	// Creates a client
	psClient, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Create the topic
	topic, err = gcp.CreateTopic(topicName, psClient)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Using topic %v to post notifications.\n", topic)

	router := mux.NewRouter()

	router.HandleFunc("/api/version", getVersionHanlder).Methods("GET")
	router.HandleFunc("/api/callpoint", postCallpointHanlder).Methods("POST")
	router.HandleFunc("/api/callpoints", getCallpointsHanlder).Methods("GET")

	log.Printf("Service: %s. Listening on port %s", appName, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func getVersionHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("[/version:GET] Requested api version.")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, fmt.Sprintf(`{"service": "%s", "version": "%s"}`, appName, appVersion))
}

func postCallpointHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoint:POST] Post a new Callpoint.")

	var n pbn.Notification

	if err := jsonpb.Unmarshal(r.Body, &n); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, fmt.Sprintf(`{"service": "%s", "version": "%s"}`, appName, appVersion))
}

func getCallpointsHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoints:GET] Requested get all callpoints")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"callpoints": {} }`)
}
