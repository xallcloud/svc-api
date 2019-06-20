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
	log.Println("[/assignments/callpoint:GET] Requested all devices associated to a callpoint")

	params := mux.Vars(r)
	cpID := params["cpID"]

	log.Println("[getEventsByCallpoint] parameter cpID:", cpID)

	if cpID == "" {
		processError(nil, w, http.StatusInternalServerError, "ERROR", "Invalid cpID!")
		return
	}

	log.Println("[getEventsByCallpoint] Create Context.")

	ctx := context.Background()

	events, err := gcp.EventsGetByCpID(ctx, dsClient, cpID)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list assignments!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	gcp.EventsToJSON(w, events)
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	log.Println("[/events:GET] Requested get all events")

	log.Println("[getEvents] Create Context.")

	ctx := context.Background()

	events, err := gcp.EventsListAll(ctx, dsClient)
	if err != nil {
		processError(err, w, http.StatusInternalServerError, "ERROR", "Could not list events!")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	gcp.EventsToJSON(w, events)
}
