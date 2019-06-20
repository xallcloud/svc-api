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

////////////////////////////////////////////////////////////////////////////////////////////////
/// Devices
////////////////////////////////////////////////////////////////////////////////////////////////

func postDevice(w http.ResponseWriter, r *http.Request) {
	log.Println("[/device:POST] Post a new Device.")

	log.Println("[postDevice] decode JSON")

	var dv pbt.Device

	if err := jsonpb.Unmarshal(r.Body, &dv); err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Bad Request! Unable to decode JSON")
		return
	}

	log.Println("[postDevice] validate JSON")

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
		Settings:    CommonSettingsToJSON(dv.Settings),
		RawRequest:  string(body),
	}

	//log.Println("[postDevice] dsDv.Settings", dsDv.Settings)

	log.Println("[postDevice] Saving message to datastore")

	ctx := context.Background()

	key, err := DeviceAdd(ctx, dsClient, dsDv)
	if err != nil && key == nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not save device to datastore!")
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

// CommonSettingsToJSON converts all the settings into a string.
func CommonSettingsToJSON(sts []*pbt.CommonSetting) string {
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
