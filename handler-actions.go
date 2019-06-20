package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	dst "github.com/xallcloud/api/datastore"
	pbt "github.com/xallcloud/api/proto"
	gcp "github.com/xallcloud/gcp"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Actions
////////////////////////////////////////////////////////////////////////////////////////////////

func postAction(w http.ResponseWriter, r *http.Request) {
	log.Println("[/action:POST] Post a new Callpoint.")

	log.Println("[postAction] decode JSON")

	var ac pbt.Action

	if err := jsonpb.Unmarshal(r.Body, &ac); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	log.Println("[postAction] validate JSON")

	if ac.AcID == "" || ac.CpID == "" || ac.Action == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Mandatory field(s) missing!")
		return
	}

	log.Println("[postAction] Encode command back to JSON.")

	ma := jsonpb.Marshaler{}
	body, err := ma.MarshalToString(&ac)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Unable to encode proto data to JSON!")
		return
	}

	if len(body) == 0 {
		processError(err, w, http.StatusBadRequest, "ERROR", "Encoding proto data to JSON: empty raw body!")
		return
	}

	dsAc := &dst.Action{
		AcID:        ac.AcID,
		CpID:        ac.CpID,
		Action:      ac.Action,
		Description: ac.Description,
		Created:     time.Now(),
		RawRequest:  string(body),
	}

	//log.Println("[postDevice] dsAc.Settings", dsAc.Settings)

	log.Println("[postDevice] Saving message to datastore")

	ctx := context.Background()

	key, err := gcp.ActionAdd(ctx, dsClient, dsAc)
	if err != nil && key == nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not save action to datastore!")
		return
	}
	ac.KeyID = key.ID

	exists := (err != nil && key != nil)

	if !exists {
		log.Printf("[datastore] stored using key: %d | %s", key.ID, ac.AcID)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf("[datastore] duplicate acID. Was stored using key: %d | %s", key.ID, ac.AcID)
		w.WriteHeader(http.StatusConflict)
	}

	//POST message to PUB/SUB
	if err := PublishAction(ctx, psClient, &ac); err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Unable to publish action")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pbt.Action{AcID: ac.AcID, KeyID: ac.KeyID})
}

func getActions(w http.ResponseWriter, r *http.Request) {
	log.Println("[/actions:GET] Requested get all actions")

	log.Println("[getActions] Create Context.")

	ctx := context.Background()

	dvs, err := gcp.ActionsListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list actions!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	gcp.ActionsToJSON(w, dvs)
}
