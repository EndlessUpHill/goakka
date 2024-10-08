package core

import (
	"sync"
	"testing"
)

// Test suite for InMemorybroker
func TestInMemorybroker(t *testing.T) {

	t.Run("TestSingleSubscriptionAndPublish", func(t *testing.T) {
		// Arrange
		received := false
		var mockwg sync.WaitGroup
		mockwg.Add(1)
		mockActor := NewBasicActor("mock-actor", func(msg *ActorResult) *ActorResult {
			received = true
			mockwg.Done()
			return msg
		})

		mockActor.Start()
		defer mockActor.Stop()

		broker := NewInMemoryBroker()
		broker.Subscribe("test-topic", mockActor)

		// Act
		broker.Publish("test-topic", "test message")

		// wait for message delivery
		mockwg.Wait()

		// Assert
		if !received {
			t.Errorf("expected message to be received by subscriber")
		}
	})

	t.Run("TestMultipleSubscriptionsAndPublish", func(t *testing.T) {
		// Arrange
		receivedMessages := make(map[string]bool)
		mu := &sync.Mutex{}

		var mockwg sync.WaitGroup

		mockwg.Add(2)

		mockActor1 := NewBasicActor("mock-actor1", func(msg *ActorResult) *ActorResult {
			mu.Lock()
			receivedMessages["mock-actor1"] = true
			mu.Unlock()
			mockwg.Done()
			return msg
		})

		mockActor1.Start()
		defer mockActor1.Stop()

		mockActor2 := NewBasicActor("mock-actor2", func(msg *ActorResult) *ActorResult {
			mu.Lock()
			receivedMessages["mock-actor2"] = true
			mu.Unlock()
			mockwg.Done()
			return msg
		})

		mockActor2.Start()
		defer mockActor2.Stop()

		broker := NewInMemoryBroker()
		broker.Subscribe("test-topic", mockActor1)
		broker.Subscribe("test-topic", mockActor2)

		// Act
		broker.Publish("test-topic", "test message")

		// wait for message delivery
		mockwg.Wait()

		// Assert
		mu.Lock()
		defer mu.Unlock()
		if !receivedMessages["mock-actor1"] {
			t.Errorf("expected message to be received by actor1")
		}
		if !receivedMessages["mock-actor2"] {
			t.Errorf("expected message to be received by actor2")
		}
	})

	// t.Run("TestConcurrentPublishAndSubscribe", func(t *testing.T) {
	// 	// Arrange
	// 	received := false

	// 	var mockwg sync.WaitGroup
	// 	mockwg.Add(1)
	// 	mockActor := NewBasicActor("mock-actor", func(msg *ActorResult) *ActorResult {
	// 			received = true
	// 			mockwg.Done()
	// 			return msg
	// 		},)

	// 	mockActor.Start()
	// 	defer mockActor.Stop()

	// 	broker := NewInMemorybroker()

	// 	// Act
	// 	var wg sync.WaitGroup
	// 	wg.Add(2)

	// 	go func() {
	// 		defer wg.Done()
	// 		broker.Subscribe("test-topic", mockActor)
	// 	}()

	// 	go func() {
	// 		defer wg.Done()
	// 		broker.Publish("test-topic", "test message")
	// 	}()

	// 	wg.Wait()

	// 	mockwg.Wait()
	// 	// Assert
	// 	if !received {
	// 		t.Errorf("expected message to be received by subscriber in concurrent scenario")
	// 	}
	// })

	t.Run("TestStressPublish", func(t *testing.T) {
		// Arrange
		receivedCount := 0
		mu := &sync.Mutex{}
		var mockwg sync.WaitGroup

		messages := 1000

		mockwg.Add(messages)
		mockActor := NewBasicActorWithMailboxSize("mock-actor2", messages,
			func(msg *ActorResult) *ActorResult {
				mu.Lock()
				receivedCount++
				mu.Unlock()
				mockwg.Done()
				return msg
			})

		mockActor.Start()
		defer mockActor.Stop()

		broker := NewInMemoryBroker()
		broker.Subscribe("test-topic", mockActor)

		// Act
		for i := 0; i < messages; i++ {
			go broker.Publish("test-topic", "message")
		}

		// wait for message delivery

		mockwg.Wait()

		// Assert
		mu.Lock()
		defer mu.Unlock()
		if receivedCount != messages {
			t.Errorf("expected to receive 1000 messages, got %d", receivedCount)
		}
	})

}
