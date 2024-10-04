package core

import "sync"

const DefaultRegistrySize = 100

// ActorRegistry holds references to all actors
type ActorRegistry struct {
    sync.RWMutex
    actors map[string]Actor
}

// NewActorRegistry creates a new instance of ActorRegistry with an optional size
func NewActorRegistry(size ...int) *ActorRegistry {
	// If size is provided, use it, otherwise default to DefaultRegistrySize
	registrySize := DefaultRegistrySize
	if len(size) > 0 {
		registrySize = size[0]
	}

	return &ActorRegistry{
		actors: make(map[string]Actor, registrySize),
	}
}

// RegisterActor adds an actor to the registry
func (ar *ActorRegistry) RegisterActor(actor Actor) {
    ar.Lock()
    defer ar.Unlock()
    ar.actors[actor.GetID()] = actor
}

// GetActor retrieves an actor by ID
func (ar *ActorRegistry) GetActor(id string) (Actor, bool) {
    ar.RLock()
    defer ar.RUnlock()
    actor, exists := ar.actors[id]
    return actor, exists
}
