package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"

	dst "github.com/xallcloud/api/datastore"
	pbt "github.com/xallcloud/api/proto"
	_ "github.com/xallcloud/gcp"
)

const (
	appName    = "svc-api"
	appVersion = "0.0.1-alfa004"
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

	log.Println("[postCallpointHanlder] decode JSON")

	var cp pbt.Callpoint

	if err := jsonpb.Unmarshal(r.Body, &cp); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	log.Println("[postCallpointHanlder] validate JSON")

	//For now, only accept "Notify" Commands
	if cp.CpID == "" && cp.AbsAddress == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Mandatory field(s) missing!")
		return
	}

	log.Println("[postCallpointHanlder] Encode command back to JSON.")

	ma := jsonpb.Marshaler{}
	body, err := ma.MarshalToString(&cp)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Unable to encode proto data to JSON!")
		return
	}

	if len(body) == 0 {
		processError(err, w, http.StatusBadRequest, "ERROR", "Encoding proto data to JSON: empty raw body!")
		return
	}

	dsCp := &dst.Callpoint{
		CpID:        cp.CpID,
		Created:     time.Now(),
		AbsAddress:  cp.AbsAddress,
		Label:       cp.Label,
		Description: cp.Description,
		Type:        cp.Type,
		Priority:    cp.Priority,
		Icon:        cp.Icon,
		RawRequest:  string(body),
	}

	log.Println("[postCallpointHanlder] Saving message to datastore")

	ctx := context.Background()

	key, err := CallpointAdd(ctx, dsClient, dsCp)
	if err != nil && key == nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not save callpoint to datastore!")
		return
	}
	cp.KeyID = key.ID

	exists := (err != nil && key != nil)

	if !exists {
		//PubNotification(ctx, psClient, &n)

		log.Printf("[datastore] stored using key: %d | %s", key.ID, cp.CpID)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf("[datastore] duplicate eventID. Was stored using key: %d | %s", key.ID, cp.CpID)
		w.WriteHeader(http.StatusConflict)
	}

	// final OK return
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pbt.Callpoint{CpID: cp.CpID, KeyID: cp.KeyID})
}

func getCallpointsHanlder(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoints:GET] Requested get all callpoints")

	log.Println("[getCallpointsHanlder] Create Context.")

	ctx := context.Background()

	cps, err := CallpointsListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list callpoints!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	CallpointsToJSON(w, cps)
}

func processError(e error, w http.ResponseWriter, httpCode int, status string, detail string) {
	log.Println(e)

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"status\":\"%s\", \"detail\":\"%s\"}", status, detail)
}
