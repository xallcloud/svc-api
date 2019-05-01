package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/datastore"

	dst "github.com/xallcloud/api/datastore"
)

//AddCallpoint method that
func AddCallpoint(ctx context.Context, client *datastore.Client, cp *dst.Callpoint) (*datastore.Key, error) {

	// first check if there already exists this Callpoint ID:
	cps, err := CallpointGetByCpID(ctx, client, cp.CpID)
	if err != nil {
		return nil, err
	}

	// if has already the value, return key and error
	if len(cps) > 0 {
		return &datastore.Key{ID: cps[0].ID, Kind: dst.KindCallpoints}, fmt.Errorf("cpId allready exists. %d", cps[0].ID)
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

// CallpointGetByCpID will return the list of callpoints with the same cpId
func CallpointGetByCpID(ctx context.Context, client *datastore.Client, cpID string) ([]*dst.Callpoint, error) {
	// Create a query to fetch all Task entities, ordered by "created".

	log.Println("[CallpointGetByCpID] will filter by:", cpID)

	var callpoints []*dst.Callpoint
	query := datastore.NewQuery(dst.KindCallpoints).
		Filter("cpId=", cpID).
		Order("created")

	log.Println("[CallpointGetByCpID] will perform query")

	keys, err := client.GetAll(ctx, query, &callpoints)
	if err != nil {
		return nil, err
	}

	log.Println("[CallpointGetByCpID] Total keys returned", len(keys))

	// Set the id field on each Task from the corresponding key.
	for i, key := range keys {
		callpoints[i].ID = key.ID
	}

	return callpoints, nil
}
