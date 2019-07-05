package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xallcloud/gcp"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Events
////////////////////////////////////////////////////////////////////////////////////////////////

func getEventsByCallpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("[/events/callpoint:GET] Requested all events associated to a callpoint")

	params := mux.Vars(r)
	cpID := params["cpID"]

	log.Println("[getEventsByCallpoint] parameter cpID:", cpID)

	if cpID == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Invalid cpID!")
		return
	}

	log.Println("[getEventsByCallpoint] Create Context.")

	ctx := context.Background()

	events, err := gcp.EventsGetByCpID(ctx, dsClient, cpID)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not list assignments!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	gcp.EventsToJSON(w, events)
}

func getEventsByAction(w http.ResponseWriter, r *http.Request) {
	log.Println("[/events/action:GET] Requested all events associated to a action/notification")

	params := mux.Vars(r)
	acID := params["acID"]

	log.Println("[getEventsByCallpoint] parameter cpID:", acID)

	if acID == "" {
		processError(nil, w, http.StatusBadRequest, "ERROR", "Invalid acID!")
		return
	}

	log.Println("[getEventsByCallpoint] Create Context.")

	ctx := context.Background()

	events, err := gcp.EventsGetByAcID(ctx, dsClient, acID)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not list events by action!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	setupResponse(w)
	gcp.EventsToJSON(w, events)
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("[/events:GET] Requested get all events")

	log.Println("[getEvents] Create Context.")

	ctx := context.Background()

	events, err := gcp.EventsListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusBadRequest, "ERROR", "Could not list events!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	gcp.EventsToJSON(w, events)
}
