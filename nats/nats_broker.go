package nats

import (
	"fmt"
	"github.com/EndlessUpHill/goakka/core"
	"github.com/nats-io/nats.go"
	"log"
)

// NatsBroker is an implementation of the MessageBroker interface using NATS Pub/Sub
type NatsBroker struct {
	conn *nats.Conn
}

// NewNatsBroker creates a new NatsBroker instance
func NewNatsBroker(natsURL string) *NatsBroker {
	// Connect to the NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	return &NatsBroker{
		conn: nc,
	}
}

// Publish sends a message to a NATS Pub/Sub topic
func (b *NatsBroker) Publish(topic string, msg interface{}) error {
	// Convert the message to a string (could use JSON or another serialization method)
	message := fmt.Sprintf("%v", msg)
	return b.conn.Publish(topic, []byte(message))
}

// Subscribe subscribes an actor to a NATS Pub/Sub topic
func (b *NatsBroker) Subscribe(topic string, actor *core.BasicActor) error {
	// Subscribe to the topic and process incoming messages
	_, err := b.conn.Subscribe(topic, func(m *nats.Msg) {
		// Pass the message payload to the actor
		actor.SendMessage(string(m.Data))
	})
	if err != nil {
		return fmt.Errorf("error subscribing to topic %s: %v", topic, err)
	}

	return nil
}
