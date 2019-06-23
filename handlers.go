package main

import (
	"fmt"
	"log"
	"net/http"
)

//processError generic error processing method to fill in default HTTP content
func processError(e error, w http.ResponseWriter, httpCode int, status string, detail string) {
	log.Println(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, `{"status":"%s", "description":"%s", "fullError":"%s"}`, status, detail, e.Error())
}

// getVersion get version
func getVersion(w http.ResponseWriter, r *http.Request) {
	log.Println("[/version:GET] Requested api version. " + appVersion)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf(`{"service": "%s", "version": "%s"}`, appName, appVersion))
}
