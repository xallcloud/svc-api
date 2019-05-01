package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/mux"
	dst "github.com/xallcloud/api/datastore"
	pbt "github.com/xallcloud/api/proto"
)

func processError(e error, w http.ResponseWriter, httpCode int, status string, detail string) {
	log.Println(e)

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"%s", "description":"%s", "fullError":"%s"}`, status, detail, e.Error())
}

////////////////////////////////////////////////////////////////////////////////////////////////
/// Version
////////////////////////////////////////////////////////////////////////////////////////////////

func getVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("[/version:GET] Requested api version.")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, fmt.Sprintf(`{"service": "%s", "version": "%s"}`, appName, appVersion))
}

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

	//For now, only accept "Notify" Commands
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pbt.Callpoint{CpID: cp.CpID, KeyID: cp.KeyID})
}

func getCallpoints(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoints:GET] Requested get all callpoints")

	log.Println("[getCallpoints] Create Context.")

	ctx := context.Background()

	cps, err := CallpointsListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list callpoints!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	CallpointsToJSON(w, cps)
}

func deleteCallpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("[/callpoint:DELETE] Requested delete a callpoint based on Key")

	params := mux.Vars(r)
	i := params["id"]

	log.Println("[deleteCallpoint] parameter id:", i)

	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not convert parameter ID to a proper number!")
		return
	}

	log.Println("[deleteCallpoint] Create Context.")

	ctx := context.Background()

	err = CallpointDelete(ctx, dsClient, id)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not delete callpoint!")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain")
}

////////////////////////////////////////////////////////////////////////////////////////////////
/// Devices
////////////////////////////////////////////////////////////////////////////////////////////////

func postDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("[/device:POST] Post a new Callpoint.")

	log.Println("[postDevice] decode JSON")

	var dv pbt.Device

	if err := jsonpb.Unmarshal(r.Body, &dv); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	log.Println("[postDevice] validate JSON")

	//For now, only accept "Notify" Commands
	if dv.DvID == "" || dv.Category == "" || dv.Destination == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Mandatory field(s) missing!")
		return
	}

	log.Println("[postDevice] Encode command back to JSON.")

	ma := jsonpb.Marshaler{}
	body, err := ma.MarshalToString(&dv)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Unable to encode proto data to JSON!")
		return
	}

	if len(body) == 0 {
		processError(err, w, http.StatusBadRequest, "ERROR", "Encoding proto data to JSON: empty raw body!")
		return
	}

	dsDv := &dst.Device{
		DvID:        dv.DvID,
		Created:     time.Now(),
		Label:       dv.Label,
		Description: dv.Description,
		Type:        dv.Type,
		Priority:    dv.Priority,
		Icon:        dv.Icon,
		IsTwoWay:    dv.IsTwoWay,
		Category:    dv.Category,
		Destination: dv.Destination,
		Settings:    DeviceSettingsToJSON(dv.Settings),
		RawRequest:  string(body),
	}

	//log.Println("[postDevice] dsDv.Settings", dsDv.Settings)

	log.Println("[postDevice] Saving message to datastore")

	ctx := context.Background()

	key, err := DeviceAdd(ctx, dsClient, dsDv)
	if err != nil && key == nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not save callpoint to datastore!")
		return
	}
	dv.KeyID = key.ID

	exists := (err != nil && key != nil)

	if !exists {
		log.Printf("[datastore] stored using key: %d | %s", key.ID, dv.DvID)
		w.WriteHeader(http.StatusCreated)
	} else {
		log.Printf("[datastore] duplicate dvID. Was stored using key: %d | %s", key.ID, dv.DvID)
		w.WriteHeader(http.StatusConflict)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pbt.Device{DvID: dv.DvID, KeyID: dv.KeyID})
}

// DeviceSettingsToJSON converts all the settings into a string.
func DeviceSettingsToJSON(sts []*pbt.DeviceSetting) string {
	var term = ""
	var r = "["
	for _, s := range sts {
		r = r + fmt.Sprintf(`%s{"name":"%s","value":"%s","type":"%s"}`,
			term,
			s.Name,
			s.Value,
			s.Type,
		)
		term = ","
	}
	r = r + "]"

	return r
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	log.Println("[/devices:GET] Requested get all devices")

	log.Println("[getDevices] Create Context.")

	ctx := context.Background()

	dvs, err := DevicesListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list devices!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	DevicesToJSON(w, dvs)
}

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("[/device:DELETE] Requested delete a device based on Key")

	params := mux.Vars(r)
	i := params["id"]

	log.Println("[deleteDevice] parameter id:", i)

	id, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not convert parameter ID to a proper number!")
		return
	}

	log.Println("[deleteDevice] Create Context.")

	ctx := context.Background()

	err = DeviceDelete(ctx, dsClient, id)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not delete device!")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain")
}

////////////////////////////////////////////////////////////////////////////////////////////////
/// Assignments
////////////////////////////////////////////////////////////////////////////////////////////////

func getAssignmentsByCallpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("[/assignments/callpoint:GET] Requested all devices associated to a callpoint")

	params := mux.Vars(r)
	cpID := params["cpID"]

	log.Println("[getAssignmentsByCallpoint] parameter cpID:", cpID)

	if cpID == "" {
		processError(nil, w, http.StatusInternalServerError, "ERROR", "Invalid cpID!")
		return
	}

	log.Println("[getAssignmentsByCallpoint] Create Context.")

	ctx := context.Background()

	asgns, err := AssignmentsByCpID(ctx, dsClient, cpID)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list assignments!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	AssignmentsToJSON(w, asgns)
}
