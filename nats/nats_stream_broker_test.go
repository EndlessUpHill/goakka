package nats_test

import (
	"testing"
	"time"

	"github.com/EndlessUpHill/goakka/core"
	coreNats "github.com/EndlessUpHill/goakka/nats"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNATSJetStreamPubSub(t *testing.T) {
	// Connect to the NATS server (make sure NATS is running on localhost:4222)
	natsPubSub, err := coreNats.NewNATSJetStreamPubSub(nats.DefaultURL, "test-stream", "test-subject", "test-consumer")
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}

	testMessage := "Hello, NATS JetStream!"

	// Create an actor
	actor := core.NewBasicActor("test-actor", func(res *core.ActorResult) *core.ActorResult {
		// Check if the message is the expected one
		assert.Equal(t, testMessage, res.Message)
		return &core.ActorResult{}
	})
	actor.Start()

	// Subscribe the actor to the NATS JetStream
	natsPubSub.Subscribe(actor)

	// Publish a message to NATS JetStream
	natsPubSub.Publish(testMessage)

	// Wait a bit for the message to be processed
	time.Sleep(100 * time.Millisecond)

	// Clean up: remove the stream from JetStream
	
	error := natsPubSub.DeleteStream("test-stream")
	
	assert.NoError(t, error)
}
