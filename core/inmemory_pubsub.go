package core

import "sync"

type PubSub interface {
	Publish(topic string, msg interface{})
	Subscribe(topic string, actor *BasicActor)
}

type InMemoryPubSub struct {
	subscribers map[string][]*BasicActor
	mu          sync.RWMutex
}

func NewInMemoryPubSub() *InMemoryPubSub {
	return &InMemoryPubSub{
		subscribers: make(map[string][]*BasicActor),
	}
}

func (ps *InMemoryPubSub) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if actors, ok := ps.subscribers[topic]; ok {
		for _, actor := range actors {
			actor.SendMessage(msg)
		}
	}
}

func (ps *InMemoryPubSub) Subscribe(topic string, actor *BasicActor) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subscribers[topic] = append(ps.subscribers[topic], actor)
}
