package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"text/tabwriter"
	"time"

	"cloud.google.com/go/datastore"

	dst "github.com/xallcloud/api/datastore"
)

////////////////////////////////////////////////////////////////////////////////////////////////
/// Callpoints
////////////////////////////////////////////////////////////////////////////////////////////////

//CallpointAdd method that
func CallpointAdd(ctx context.Context, client *datastore.Client, cp *dst.Callpoint) (*datastore.Key, error) {

	// first check if there already exists this Callpoint ID:
	cps, err := CallpointGetByCpID(ctx, client, cp.CpID)
	if err != nil {
		return nil, err
	}

	// if has already the value, return key and error
	if len(cps) > 0 {
		return &datastore.Key{ID: cps[0].ID, Kind: dst.KindCallpoints}, fmt.Errorf("cpID allready exists. %d", cps[0].ID)
	}

	// copy to new record
	n := &dst.Callpoint{
		CpID:        cp.CpID,
		Created:     time.Now(),
		AbsAddress:  cp.AbsAddress,
		Label:       cp.Label,
		Description: cp.Description,
		Type:        cp.Type,
		Priority:    cp.Priority,
		Icon:        cp.Icon,
		RawRequest:  cp.RawRequest,
	}

	// do the insert
	key := datastore.IncompleteKey(dst.KindCallpoints, nil)
	return client.Put(ctx, key, n)
}

// CallpointGetByCpID will return the list of callpoints with the same cpID
func CallpointGetByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Callpoint, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[CallpointGetByCpID] will filter by cpID:", cpID)

	var callpoints []*dst.Callpoint
	query := datastore.NewQuery(dst.KindCallpoints).
		Filter("cpID =", cpID)

	log.Println("[CallpointGetByCpID] will perform query")

	keys, err := client.GetAll(ctx, query, &callpoints)
	if err != nil {
		return nil, err
	}

	log.Println("[CallpointGetByCpID] Total keys returned", len(keys))

	// Set the ID field on each Callpoint from the corresponding key.
	for i, key := range keys {
		callpoints[i].ID = key.ID
	}

	return callpoints, nil
}

// CallpointsListAll returns all the tasks in ascending order of creation time.
func CallpointsListAll(ctx context.Context, client *datastore.Client) ([]*dst.Callpoint, error) {
	var cps []*dst.Callpoint

	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery(dst.KindCallpoints).Order("created")
	keys, err := client.GetAll(ctx, query, &cps)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Callpoint from the corresponding DataStore key.
	for i, key := range keys {
		cps[i].ID = key.ID
	}

	return cps, nil
}

// CallpointDelete will delete a callpoint from the datastore
func CallpointDelete(ctx context.Context, client *datastore.Client, cpKeyID int64) error {
	return client.Delete(ctx, datastore.IDKey(dst.KindCallpoints, cpKeyID, nil))
}

// CallpointsToJSON prints the callpoints into JSON to the given writer.
func CallpointsToJSON(w io.Writer, cps []*dst.Callpoint) {
	const line = `%s
	{
		"ID": %d,
		"cpID": "%s",
		"created": "%v",
		"absAddress": "%s",
		"label": "%s",
		"description": "%s",
		"type": %d,
		"priority": %d,
		"icon": "%s",
		"rawRequest": %s
	}`

	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.

	var term = ""
	var rawRequest string
	fmt.Fprintf(tw, "[\n")
	for _, c := range cps {
		rawRequest = strings.TrimSpace(c.RawRequest)

		//log.Println("[JSONNotifications] parameter raw:", raw)

		if rawRequest == "" {
			rawRequest = "null"
		}

		fmt.Fprintf(tw, line, term,
			c.ID,
			c.CpID,
			c.Created,
			c.AbsAddress,
			c.Label,
			c.Description,
			c.Type,
			c.Priority,
			c.Icon,
			rawRequest,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

////////////////////////////////////////////////////////////////////////////////////////////////
/// Devices
////////////////////////////////////////////////////////////////////////////////////////////////

//DeviceAdd method that
func DeviceAdd(ctx context.Context, client *datastore.Client, dv *dst.Device) (*datastore.Key, error) {

	// first check if there already exists this Callpoint ID:
	cps, err := DeviceGetByDvID(ctx, client, dv.DvID)
	if err != nil {
		return nil, err
	}

	// if has already the value, return key and error
	if len(cps) > 0 {
		return &datastore.Key{ID: cps[0].ID, Kind: dst.KindCallpoints}, fmt.Errorf("dvID allready exists. %d", cps[0].ID)
	}

	// copy to new record
	n := &dst.Device{
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
		Settings:    dv.Settings,
		RawRequest:  dv.RawRequest,
	}

	// do the insert
	key := datastore.IncompleteKey(dst.KindDevices, nil)
	return client.Put(ctx, key, n)
}

// DeviceGetByDvID will return the list of devices with the same dvID
func DeviceGetByDvID(ctx context.Context, client *datastore.Client, dvID string) ([]*dst.Device, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[DeviceGetByDvID] will filter by cpID:", dvID)

	var devices []*dst.Device
	query := datastore.NewQuery(dst.KindDevices).
		Filter("dvID =", dvID)

	log.Println("[DeviceGetByDvID] will perform query")

	keys, err := client.GetAll(ctx, query, &devices)
	if err != nil {
		return nil, err
	}

	log.Println("[DeviceGetByDvID] Total keys returned", len(keys))

	// Set the ID field on each Callpoint from the corresponding key.
	for i, key := range keys {
		devices[i].ID = key.ID
	}

	return devices, nil
}

// DevicesListAll returns all the devices in ascending order of creation time.
func DevicesListAll(ctx context.Context, client *datastore.Client) ([]*dst.Device, error) {
	var dvs []*dst.Device

	// Create a query to fetch all Devices entities, ordered by "created".
	query := datastore.NewQuery(dst.KindDevices).Order("created")
	keys, err := client.GetAll(ctx, query, &dvs)
	if err != nil {
		return nil, err
	}

	// Set the id field on each Devices from the corresponding DataStore key.
	for i, key := range keys {
		dvs[i].ID = key.ID
	}

	return dvs, nil
}

// DeviceDelete will delete a device from the datastore
func DeviceDelete(ctx context.Context, client *datastore.Client, dvKeyID int64) error {
	return client.Delete(ctx, datastore.IDKey(dst.KindDevices, dvKeyID, nil))
}

// DevicesToJSON prints the callpoints into JSON to the given writer.
func DevicesToJSON(w io.Writer, dvs []*dst.Device) {
	const line = `%s
	{
		"ID": %d,
		"dvID": "%s",
		"created": "%v",
		"label": "%s",
		"description": "%s",
		"type": %d,
		"priority": %d,
		"isTwoWay": %s,
		"category": "%s",
		"destination": "%s",
		"icon": "%s",
		"settings": %s,
		"rawRequest": %s
	}`

	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.

	var term = ""
	var rawRequest, isTwoWayString string
	fmt.Fprintf(tw, "[\n")
	for _, d := range dvs {
		rawRequest = strings.TrimSpace(d.RawRequest)

		if d.IsTwoWay {
			isTwoWayString = "true"
		} else {
			isTwoWayString = "false"
		}

		//log.Println("[JSONNotifications] parameter raw:", raw)

		if rawRequest == "" {
			rawRequest = "null"
		}

		fmt.Fprintf(tw, line, term,
			d.ID,
			d.DvID,
			d.Created,
			d.Label,
			d.Description,
			d.Type,
			d.Priority,
			isTwoWayString,
			d.Category,
			d.Destination,
			d.Icon,
			d.Settings,
			rawRequest,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}

////////////////////////////////////////////////////////////////////////////////////////////////
/// Assignments
////////////////////////////////////////////////////////////////////////////////////////////////

// AssignmentsByCpID will return the list of assignments and its information with the same cpID
func AssignmentsByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Assignment, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[AssignmentsByCpID] will filter by cpID:", cpID)

	var assignments []*dst.Assignment
	query := datastore.NewQuery(dst.KindAssignments).
		Filter("cpID =", cpID)

	log.Println("[AssignmentsByCpID] will perform query")

	keys, err := client.GetAll(ctx, query, &assignments)
	if err != nil {
		return nil, err
	}

	log.Println("[AssignmentsByCpID] Total keys returned", len(keys))

	// Set the ID field on each Callpoint from the corresponding key.
	for i, key := range keys {
		assignments[i].ID = key.ID
	}

	return assignments, nil
}

// AssignmentsToJSON prints the assignments into JSON to the given writer.
func AssignmentsToJSON(w io.Writer, asgs []*dst.Assignment) {
	const line = `%s
	{
		"ID": %d,
		"asID": "%s",
		"created": "%v",
		"changed": "%v",
		"description": "%s",
		"cpID": "%s",
		"dvID": "%s",
		"level": %d,
		"settings": %s,
		"rawRequest": %s
	}`

	// Use a tab writer to help make results pretty.
	tw := tabwriter.NewWriter(w, 4, 4, 1, ' ', 0) // Min cell size of 8.

	var term = ""
	var rawRequest string
	fmt.Fprintf(tw, "[\n")
	for _, a := range asgs {
		rawRequest = strings.TrimSpace(a.RawRequest)

		if rawRequest == "" {
			rawRequest = "null"
		}

		fmt.Fprintf(tw, line, term,
			a.ID,
			a.DvID,
			a.Created,
			a.Changed,
			a.Description,
			a.CpID,
			a.DvID,
			a.Level,
			a.Settings,
			rawRequest,
		)
		term = ","
	}
	fmt.Fprintf(tw, "\n]")
	tw.Flush()
}
