package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test suite for BasicActor
func TestActor(t *testing.T) {

	t.Run("TestActorStartAndStop", func(t *testing.T) {
		// Arrange
		wg := &sync.WaitGroup{}
		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = func(result *ActorResult) *ActorResult {
			return &ActorResult{}
		}
		actor.SetWaitGroup(wg)

		// Act
		actor.Start()
		actor.Stop()

		// Assert
		wg.Wait() // Ensure actor has fully stopped
	})

	t.Run("TestActorReceivesMessage", func(t *testing.T) {
		// Arrange
		received := false
		receiveFunc := func(result *ActorResult) *ActorResult {
			received = true
			return &ActorResult{}
		}
		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = receiveFunc

		// Act
		actor.Start()
		actor.SendMessage("test message")
		time.Sleep(100 * time.Millisecond) // Give some time for the actor to process

		// Assert
		if !received {
			t.Errorf("expected actor to receive the message")
		}
		actor.Stop()
	})

	t.Run("TestActorHandlesContextCancellation", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = func(result *ActorResult) *ActorResult {
			return &ActorResult{}
		}
		actor.SetWaitGroup(wg)
		actor.SetContext(ctx)

		// Act
		actor.Start()
		cancel() // Cancel the context to trigger actor shutdown

		// Assert
		wg.Wait() // Ensure actor has fully stopped
	})

	t.Run("TestActorFailureChannel", func(t *testing.T) {
		// Arrange
		failureChannel := make(chan *ActorResult, 1)
		receiveFunc := func(result *ActorResult) *ActorResult {
			return &ActorResult{
				Error: errors.New("test failure"),
			}
		}
		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = receiveFunc
		actor.SetFailureChannel(failureChannel)

		// Act
		actor.Start()
		actor.SendMessage("trigger failure")
		time.Sleep(100 * time.Millisecond) // Allow some time for message processing

		// Assert
		select {
		case failure := <-failureChannel:
			if failure.Error == nil {
				t.Errorf("expected an error in the failure channel")
			}
		default:
			t.Errorf("expected a failure to be sent to the failure channel")
		}

		actor.Stop()
	})

	t.Run("TestActorMailboxFull", func(t *testing.T) {
		// Arrange
		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = func(result *ActorResult) *ActorResult {
			return &ActorResult{}
		}
		actor.Start()

		// Fill the mailbox
		for i := 0; i < 100; i++ {
			actor.SendMessage(i)
		}

		// Act
		actor.SendMessage("overflow message") // This message should be dropped

		// Assert
		time.Sleep(100 * time.Millisecond) // Allow some time for message processing

		actor.Stop()
		// Since the mailbox is full, the actor should drop the overflow message
		// Manually check for a logged message or ensure no panic
	})

	t.Run("TestConcurrentMessageSending", func(t *testing.T) {
		// Arrange
		receivedMessages := make([]string, 0)
		mu := &sync.Mutex{}
		receiveFunc := func(result *ActorResult) *ActorResult {
			mu.Lock()
			receivedMessages = append(receivedMessages, result.Message.(string))
			mu.Unlock()
			return &ActorResult{}
		}

		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = receiveFunc
		actor.Start()

		// Act
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				actor.SendMessage(fmt.Sprintf("message-%d", i))
			}(i)
		}

		wg.Wait()
		time.Sleep(100 * time.Millisecond) // Give some time for the actor to process

		// Assert
		if len(receivedMessages) != 100 {
			t.Errorf("expected to receive 100 messages, got %d", len(receivedMessages))
		}
		actor.Stop()
	})

	t.Run("TestActorFailureAndRecovery", func(t *testing.T) {
		// Arrange
		failureChannel := make(chan *ActorResult, 10)
		recovered := false

		receiveFunc := func(result *ActorResult) *ActorResult {
			if result.Message == "fail" {
				return &ActorResult{Error: errors.New("intentional failure")}
			} else if result.Message == "recover" {
				recovered = true
				return &ActorResult{}
			}
			return &ActorResult{}
		}

		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = receiveFunc
		actor.SetFailureChannel(failureChannel)
		actor.Start()

		// Act
		actor.SendMessage("fail")
		time.Sleep(100 * time.Millisecond) // Allow time for failure
		actor.SendMessage("recover")
		time.Sleep(100 * time.Millisecond) // Allow time for recovery

		// Assert
		select {
		case failure := <-failureChannel:
			if failure.Error == nil {
				t.Errorf("expected an error in the failure channel")
			}
		default:
			t.Errorf("expected a failure to be sent to the failure channel")
		}

		if !recovered {
			t.Errorf("expected actor to recover from failure")
		}

		actor.Stop()
	})

	t.Run("TestParentChildContextCancellation", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithCancel(context.Background())
		childWg := &sync.WaitGroup{}
		parentWg := &sync.WaitGroup{}

		childActor := NewBasicActor("child-actor")
		childActor.SetWaitGroup(childWg)
		childActor.SetContext(ctx)
		childActor.ReceiveFunc = func(result *ActorResult) *ActorResult {
			return &ActorResult{}
		}

		parentActor := NewBasicActor("parent-actor")
		parentActor.ReceiveFunc = func(result *ActorResult) *ActorResult {
			if result.Message == "cancel" {
				cancel() // Cancel the context, causing the child to stop
			}
			return &ActorResult{}
		}
		parentActor.SetWaitGroup(parentWg)
		parentActor.SetContext(ctx)

		// Start both actors
		childActor.Start()
		parentActor.Start()

		// Act
		parentActor.SendMessage("cancel")
		time.Sleep(100 * time.Millisecond) // Give time for context cancellation

		// Assert
		parentWg.Wait()
		childWg.Wait() // Ensure both parent and child actors have stopped

		if len(parentActor.mailbox) != 0 || len(childActor.mailbox) != 0 {
			t.Errorf("expected both actors to stop and have empty mailboxes")
		}
	})

	t.Run("TestInterActorCommunication", func(t *testing.T) {
		// Arrange
		childReceived := false
		childReceiveFunc := func(result *ActorResult) *ActorResult {
			if result.Message == "ping" {
				childReceived = true
			}
			return &ActorResult{}
		}

		parentActor := NewBasicActor("parent-actor")
		childActor := NewBasicActor("child-actor")
		childActor.ReceiveFunc = childReceiveFunc

		parentActor.Start()
		childActor.Start()

		// Act: Parent sends a message to the child
		childActor.SendMessage("ping")
		time.Sleep(100 * time.Millisecond) // Give time for message passing

		// Assert
		if !childReceived {
			t.Errorf("expected child actor to receive 'ping' message")
		}

		parentActor.Stop()
		childActor.Stop()
	})

	t.Run("TestStressMailboxCapacity", func(t *testing.T) {
		// Arrange
		actor := NewBasicActor("test-actor")
		actor.ReceiveFunc = func(result *ActorResult) *ActorResult {
			return &ActorResult{}
		}
		actor.Start()

		// Act: Send more messages than the mailbox can handle
		for i := 0; i < 150; i++ {
			actor.SendMessage(fmt.Sprintf("message-%d", i))
		}

		time.Sleep(100 * time.Millisecond) // Give time for message processing

		// Assert: Ensure no panic or crash, and mailbox is still functioning
		actor.SendMessage("final-message")
		time.Sleep(50 * time.Millisecond) // Ensure final message is processed

		actor.Stop()
	})

}
