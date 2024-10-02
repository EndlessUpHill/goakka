package core

import (
	"context"
	"fmt"
	"sync"
)

type Supervisor struct {
	actors         []*BasicActor
	subSupervisors []*Supervisor
	stop           chan struct{}
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	broker         MessageBroker
}

// NewSupervisor creates a new supervisor with an optional timeout
func NewSupervisor(ctx context.Context, broker MessageBroker) *Supervisor {
	ctx, cancel := context.WithCancel(ctx)
	if broker == nil {
		broker = NewInMemoryBroker()
	}
	return &Supervisor{
		actors:         make([]*BasicActor, 0),
		subSupervisors: make([]*Supervisor, 0),
		stop:           make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
		broker:         broker,
	}
}

// SuperviseActor adds an actor to the supervisor and starts it
func (s *Supervisor) SuperviseActor(actor *BasicActor) {
	fmt.Println("Supervisor supervising actor...")
	actor.setWaitGroup(&s.wg)
	actor.setContext(s.ctx)
	s.actors = append(s.actors, actor)
	actor.Start()
}

// SuperviseSupervisor adds a nested supervisor (creating a hierarchy)
func (s *Supervisor) SuperviseSupervisor(subSupervisor *Supervisor) {
	fmt.Println("Supervisor supervising sub-supervisor...")
	s.subSupervisors = append(s.subSupervisors, subSupervisor)
}

// Stop gracefully stops all actors and nested supervisors
func (s *Supervisor) Stop() {
	fmt.Println("Supervisor stopping all actors and sub-supervisors...")

	s.cancel()

	// Signal all actors to stop
	for _, actor := range s.actors {
		actor.Stop()
	}

	// Signal all sub-supervisors to stop
	for _, subSupervisor := range s.subSupervisors {
		subSupervisor.Stop()
	}

	// Wait for actors and supervisors to finish
	fmt.Println("Waiting for actors to stop...")
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All actors and supervisors have stopped.")
	case <-s.ctx.Done():
		fmt.Println("Context canceled or timeout reached before all actors could stop.")
	}

	close(s.stop)
}

// Wait blocks until the supervisor is stopped
func (s *Supervisor) Wait() {
	<-s.stop
	fmt.Println("Supervisor shutdown complete.")
}
