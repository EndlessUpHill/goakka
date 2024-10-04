package nats

import (
	"log"
	"time"

	"github.com/EndlessUpHill/goakka/core"

	"github.com/nats-io/nats.go"
)

type NATSJetStreamPubSub struct {
	conn       *nats.Conn
	jetStream  nats.JetStreamContext
	streamName string
	subject    string
	consumer   string
}

// NewNATSJetStreamPubSub creates a new NATS JetStream-based pub/sub system
func NewNATSJetStreamPubSub(url, streamName, subject, consumer string) (*NATSJetStreamPubSub, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	jetStream, err := conn.JetStream()
	if err != nil {
		return nil, err
	}

	// Create the stream if it doesn't exist
	_, err = jetStream.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
	})
	if err != nil {
		return nil, err
	}

	return &NATSJetStreamPubSub{
		conn:       conn,
		jetStream:  jetStream,
		streamName: streamName,
		subject:    subject,
		consumer:   consumer,
	}, nil
}

// Publish a message to the NATS JetStream
func (n *NATSJetStreamPubSub) Publish(msg string) {
	_, err := n.jetStream.Publish(n.subject, []byte(msg))
	if err != nil {
		log.Printf("Error publishing to NATS JetStream: %v", err)
	}
}

// Subscribe an actor to the NATS JetStream
func (n *NATSJetStreamPubSub) Subscribe(actor *core.BasicActor) {
	go func() {
		_, err := n.jetStream.QueueSubscribe(n.subject, n.consumer, func(msg *nats.Msg) {
			actor.SendMessage(string(msg.Data))
			// Acknowledge the message after processing
			msg.Ack()
		}, nats.ManualAck())
		if err != nil {
			log.Fatalf("Error subscribing to JetStream: %v", err)
		}
	}()

	time.Sleep(1 * time.Second) // Simulate a delay to ensure message processing
}
