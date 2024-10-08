package core

import "sync"

const DefaultRegistrySize = 100

type ActorRegistry struct {
	sync.RWMutex
	actors map[string]Actor
}

func NewActorRegistry(size ...int) *ActorRegistry {
	registrySize := DefaultRegistrySize
	if len(size) > 0 {
		registrySize = size[0]
	}

	return &ActorRegistry{
		actors: make(map[string]Actor, registrySize),
	}
}

func (ar *ActorRegistry) RegisterActor(actor Actor) {
	ar.Lock()
	defer ar.Unlock()
	ar.actors[actor.GetName()] = actor
}

func (ar *ActorRegistry) GetActor(id string) (Actor, bool) {
	ar.RLock()
	defer ar.RUnlock()
	actor, exists := ar.actors[id]
	return actor, exists
}
