package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/gogo/protobuf/proto"

	pbt "github.com/xallcloud/api/proto"
)

// PublishAction will publish a action to the pubsub stream
func PublishAction(ctx context.Context, client *pubsub.Client, ac *pbt.Action) error {
	log.Printf("[PublishAction] [acID=%s] [cpID=%s] New Action: %s", ac.AcID, ac.CpID, ac.Action)

	m, err := proto.Marshal(ac)
	if err != nil {
		return fmt.Errorf("unable to serialize data. %v", err)
	}

	msg := &pubsub.Message{
		Data: m,
	}
	var mID string
	mID, err = tcPubNot.Publish(ctx, msg).Get(ctx)
	if err != nil {
		return fmt.Errorf("could not publish message. %v", err)
	}

	log.Printf("[PublishAction] [acID=%s] New Notification published. [mID=%s]", ac.AcID, mID)

	return nil
}
