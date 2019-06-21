package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	dst "github.com/xallcloud/api/datastore"
	pbt "github.com/xallcloud/api/proto"
	"github.com/xallcloud/gcp"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Assignments
////////////////////////////////////////////////////////////////////////////////////////////////

func postAssignment(w http.ResponseWriter, r *http.Request) {
	log.Println("[/assignment:POST] Post a new Assignment.")

	log.Println("[postAssignment] decode JSON")

	var asgn pbt.Assignment

	if err := jsonpb.Unmarshal(r.Body, &asgn); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	log.Println("[postAssignment] validate JSON")

	//For now, only accept "Notify" Commands
	if asgn.AsID == "" || asgn.CpID == "" || asgn.DvID == "" || asgn.Level <= 0 {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Mandatory field(s) missing!")
		return
	}

	log.Println("[postAssignment] Encode command back to JSON.")

	ma := jsonpb.Marshaler{}
	body, err := ma.MarshalToString(&asgn)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Unable to encode proto data to JSON!")
		return
	}

	if len(body) == 0 {
		processError(err, w, http.StatusBadRequest, "ERROR", "Encoding proto data to JSON: empty raw body!")
		return
	}

	dsAsgn := &dst.Assignment{
		AsID:        asgn.AsID,
		Created:     time.Now(),
		Changed:     time.Now(),
		Description: asgn.Description,
		CpID:        asgn.CpID,
		DvID:        asgn.DvID,
		Level:       asgn.Level,
		Settings:    CommonSettingsToJSON(asgn.Settings),
		RawRequest:  string(body),
	}

	//log.Println("[postAssignment] dsDv.Settings", dsDv.Settings)

	log.Println("[postAssignment] Saving message to datastore")

	ctx := context.Background()

	key, err := gcp.AssignmentAdd(ctx, dsClient, dsAsgn)
	if err != nil && key == nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not save assignment to datastore!")
		return
	}
	asgn.KeyID = key.ID

	exists := (err != nil && key != nil)

	if !exists {
		log.Printf("[postAssignment:datastore] stored using key: %d | %s", key.ID, asgn.AsID)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf("[postAssignment:datastore] duplicate dvID. Was stored using key: %d | %s", key.ID, asgn.AsID)
		w.WriteHeader(http.StatusConflict)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pbt.Assignment{AsID: asgn.AsID, KeyID: asgn.KeyID})
}

func getAssignmentsByCallpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("[/assignments/callpoint:GET] Requested all devices associated to a callpoint")

	params := mux.Vars(r)
	cpID := params["cpID"]

	log.Println("[getAssignmentsByCallpoint] parameter cpID:", cpID)

	if cpID == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Invalid cpID!")
		return
	}

	log.Println("[getAssignmentsByCallpoint] Create Context.")

	ctx := context.Background()

	asgns, err := gcp.AssignmentsByCpID(ctx, dsClient, cpID)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not list assignments!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	gcp.AssignmentsToJSON(w, asgns)
}
