package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	dst "github.com/xallcloud/api/datastore"
	pbt "github.com/xallcloud/api/proto"
	"github.com/xallcloud/gcp"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Callpoints
////////////////////////////////////////////////////////////////////////////////////////////////

func postCallpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoint:POST] Post a new Callpoint.")

	log.Println("[postCallpoint] decode JSON")

	var cp pbt.Callpoint

	if err := jsonpb.Unmarshal(r.Body, &cp); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	log.Println("[postCallpoint] validate JSON")

	if cp.CpID == "" && cp.AbsAddress == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Mandatory field(s) missing!")
		return
	}

	log.Println("[postCallpoint] Encode command back to JSON.")

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

	log.Println("[postCallpoint] Saving message to datastore")

	ctx := context.Background()

	key, err := gcp.CallpointAdd(ctx, dsClient, dsCp)
	if err != nil && key == nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not save callpoint to datastore!")
		return
	}
	cp.KeyID = key.ID

	exists := (err != nil && key != nil)

	w.Header().Set("Content-Type", "application/json")
	if !exists {
		log.Printf("[datastore] stored using key: %d | %s", key.ID, cp.CpID)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf("[datastore] duplicate eventID. Was stored using key: %d | %s", key.ID, cp.CpID)
		w.WriteHeader(http.StatusConflict)
	}
	json.NewEncoder(w).Encode(pbt.Callpoint{CpID: cp.CpID, KeyID: cp.KeyID})
}

func getCallpoints(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoints:GET] Requested get all callpoints")

	log.Println("[getCallpoints] Create Context.")

	ctx := context.Background()

	cps, err := gcp.CallpointsListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not list callpoints!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	gcp.CallpointsToJSON(w, cps)
}

func deleteCallpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoint:DELETE] Requested delete a callpoint based on Key")

	params := mux.Vars(r)
	i := params["id"]

	log.Println("[deleteCallpoint] parameter id:", i)

	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not convert parameter ID to a proper number!")
		return
	}

	log.Println("[deleteCallpoint] Create Context.")

	ctx := context.Background()

	err = gcp.CallpointDelete(ctx, dsClient, id)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not delete callpoint!")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain")
}
