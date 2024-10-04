package core

import (
	"fmt"
	"sync"
)

type MessageBroker interface {
	Publish(topic string, msg interface{}) error
	Subscribe(topic string, actor Actor) error
}

// InMemoryBroker is an in-memory implementation of the MessageBroker interface
type InMemoryBroker struct {
	subscribers map[string][]Actor
	mu          sync.RWMutex
}

// NewInMemoryBroker creates a new in-memory broker
func NewInMemoryBroker() *InMemoryBroker {
	return &InMemoryBroker{
		subscribers: make(map[string][]Actor),
	}
}

// Publish sends a message to all actors subscribed to the topic
func (b *InMemoryBroker) Publish(topic string, msg interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	actors, ok := b.subscribers[topic]
	if !ok {
		return fmt.Errorf("no subscribers for topic %s", topic)
	}

	for _, actor := range actors {
		actor.SendMessage(msg)
	}

	return nil
}

// Subscribe adds an actor to the list of subscribers for a given topic
func (b *InMemoryBroker) Subscribe(topic string, actor Actor) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], actor)
	return nil
}
