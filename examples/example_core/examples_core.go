package main

import (
	"context"
	"fmt"

	"github.com/EndlessUpHill/goakka/core"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	fmt.Println("Starting application...")

	// Create a root context with a 10-second timeout for testing
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure cancellation when the main function exits

	// Create a top-level supervisor with the root context
	supervisor := core.NewSupervisor(ctx, nil)

	// Create actors
	actor1 := core.NewBasicActor()
	actor2 := core.NewBasicActor()

	// Supervise actors with the top-level supervisor
	supervisor.SuperviseActor(actor1)
	supervisor.SuperviseActor(actor2)

	// Create a child supervisor with its own context (inherited from the parent)
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel() // Clean up the child context

	childSupervisor := core.NewSupervisor(childCtx, nil)

	// Create actors for the child supervisor
	actor3 := core.NewBasicActor()
	actor4 := core.NewBasicActor()

	// Supervise actors with the child supervisor
	childSupervisor.SuperviseActor(actor3)
	childSupervisor.SuperviseActor(actor4)

	// Supervise the child supervisor with the top-level supervisor
	supervisor.SuperviseSupervisor(childSupervisor)

	// Send test messages to all actors
	fmt.Println("Sending test messages to actors...")
	actor1.SendMessage("Message for actor 1")
	actor2.SendMessage("Message for actor 2")
	actor3.SendMessage("Message for actor 3")
	actor4.SendMessage("Message for actor 4")

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

}
