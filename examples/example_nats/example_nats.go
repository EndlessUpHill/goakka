package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EndlessUpHill/goakka/core"
)

func main() {
	fmt.Println("Starting application...")

	// Create a root context with a 10-second timeout for testing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure cancellation when the main function exits

	// Create a top-level supervisor with the root context
	supervisor := core.NewSupervisor(ctx)

	// Create actors
	// Create and register actor1
	actor1 := core.NewBasicActor("actor1", func(msg interface{}) {
		fmt.Printf("Actor1 received: %v\n", msg)
		// Send message to actor2
		if actor2Ref, exists := NatsRegistryInstance.GetActor("actor2"); exists {
			actor2Ref.SendMessage("Hello from Actor1")
		}
	})

	NatsRegistryInstance.RegisterActor(actor1)
	supervisor.SuperviseActor(actor1)

	actor2 := core.NewBasicActor("actor2", func(msg interface{}) {
		fmt.Printf("Actor2 received: %v\n", msg)
	})

	NatsRegistryInstance.RegisterActor(actor2)
	supervisor.SuperviseActor(actor2)

	// Supervise actors with the top-level supervisor
	supervisor.SuperviseActor(actor1)
	supervisor.SuperviseActor(actor2)

	// Create a child supervisor with its own context (inherited from the parent)
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel() // Clean up the child context

	childSupervisor := core.NewSupervisor(childCtx)

	actor3 := NewExampleNats("actor3")

	// Supervise actors with the child supervisor
	childSupervisor.SuperviseActor(actor3)

	// Supervise the child supervisor with the top-level supervisor
	supervisor.SuperviseSupervisor(childSupervisor)

	// Send test messages to all actors
	fmt.Println("Sending test messages to actors...")

	actor1.SendMessage("Message for actor 1")
	actor2.SendMessage("Message for actor 2")
	actor3.SendMessage("Message for actor 3")

	NatsBrokerInstance.Publish("example", "Subscribe to the example topic to receive this message.")
	// Publish a message to the broker

	NatsBrokerInstance.Subscribe("example", actor1)
	NatsBrokerInstance.Subscribe("example", actor3)
	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Waiting for termination signal...")
		<-sigs
		fmt.Println("Termination signal received. Stopping supervisor...")
		supervisor.Stop() // Gracefully stop all supervisors and actors
	}()

	// Keep the application running until the supervisor shuts down
	supervisor.Wait()

	fmt.Println("Application shutdown complete.")
	// Create a NATS pub/sub system (replace with your NATS server URL)
}
