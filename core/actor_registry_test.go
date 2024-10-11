package core

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestActorRegistry(t *testing.T) {

	t.Run("TestNewActorRegistryWithDefaultSize", func(t *testing.T) {
		// Arranges
		registry := NewActorRegistry()

		// Act & Assert
		if len(registry.actors) != 0 {
			t.Errorf("expected empty actor map, got %d actors", len(registry.actors))
		}
	})

	t.Run("TestNewActorRegistryWithCustomSize", func(t *testing.T) {
		// Arrange
		customSize := 50
		registry := NewActorRegistry(customSize)

		// Act & Assert
		if len(registry.actors) != 0 {
			t.Errorf("expected empty actor map, got %d actors", len(registry.actors))
		}
	})

	t.Run("TestRegisterActor", func(t *testing.T) {
		// Arrange
		registry := NewActorRegistry()
		mockActor := &MockActor{name: "test-actor", id: uuid.New()}

		// Act
		registry.RegisterActor(mockActor)

		// Assert
		actor, exists := registry.GetActor("test-actor")
		if !exists {
			t.Errorf("expected actor to be registered")
		}
		if actor.GetName() != "test-actor" {
			t.Errorf("expected actor name to be 'test-actor', got '%s'", actor.GetName())
		}
	})

	t.Run("TestGetActorNotFound", func(t *testing.T) {
		// Arrange
		registry := NewActorRegistry()

		// Act
		_, exists := registry.GetActor("non-existent-actor")

		// Assert
		if exists {
			t.Errorf("expected no actor to be found")
		}
	})

	t.Run("TestConcurrentRegisterAndGetActor", func(t *testing.T) {
		// Arrange
		registry := NewActorRegistry()
		mockActor := &MockActor{name: "concurrent-actor"}
		var wg sync.WaitGroup

		// Act
		wg.Add(2)

		go func() {
			defer wg.Done()
			registry.RegisterActor(mockActor)
		}()

		go func() {
			defer wg.Done()
			_, exists := registry.GetActor("concurrent-actor")
			if exists {
				t.Errorf("actor should not exist before registration completes")
			}
		}()

		wg.Wait()

		// Assert
		actor, exists := registry.GetActor("concurrent-actor")
		if !exists {
			t.Errorf("expected actor to be registered")
		}
		if actor.GetName() != "concurrent-actor" {
			t.Errorf("expected actor name to be 'concurrent-actor', got '%s'", actor.GetName())
		}
	})
}

// Mock Actor to use for testing purposes
type MockActor struct {
	name        string
	id          uuid.UUID
}

func (ma *MockActor) GetName() string {
	return ma.name
}

func (ma *MockActor) Start() {}

func (ma *MockActor) Stop() {}

func (ma *MockActor) ReceiveMessage(msg interface{}) *ActorResult {
	return &ActorResult{}
}

func (ma *MockActor) GetContext() context.Context {
	return context.Background()
}

func (ma *MockActor) SendMessage(msg interface{}) {}

func (ma *MockActor) GetID() uuid.UUID {
	return ma.id
}

func (ma *MockActor) SetWaitGroup(wg *sync.WaitGroup) {}

func (ma *MockActor) SetContext(ctx context.Context) {}

func (ma *MockActor) SetFailureChannel(failure chan *ActorResult) {}
