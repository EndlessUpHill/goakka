package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type SupervisorActorResult struct {
	Action int
	Result *ActorResult
}

type Supervisor struct {
	id                uuid.UUID // UUID for each supervisor
	child             bool
	actors            map[uuid.UUID]Actor // Changed to a map with actor name as key
	subSupervisors    map[uuid.UUID]*Supervisor
	stop              chan struct{}
	wg                sync.WaitGroup
	ctx               context.Context
	cancel            context.CancelFunc
	actorMonitor      *ActorMonitor
	supervisorMonitor *SupervisorMonitor
}

// NewSupervisor creates a new supervisor with an optional timeout
func NewSupervisor(ctx context.Context) *Supervisor {
	ctx, cancel := context.WithCancel(ctx)
	s := &Supervisor{
		id:             uuid.New(),
		actors:         make(map[uuid.UUID]Actor),
		subSupervisors: make(map[uuid.UUID]*Supervisor),
		stop:           make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
	}
	// Initialize monitors
	s.supervisorMonitor = NewSupervisorMonitor(s)
	s.actorMonitor = NewActorMonitor(s)
	return s
}

func (s *Supervisor) GetID() uuid.UUID {
	return s.id
}

// SuperviseActor adds an actor to the supervisor and starts it
func (s *Supervisor) SuperviseActor(actor Actor) {
	fmt.Println("Supervisor supervising actor...")
	actor.SetWaitGroup(&s.wg)
	actor.SetContext(s.ctx)
	actor.SetFailureChannel(s.actorMonitor.GetInboundChannel())
	s.actors[actor.GetID()] = actor
	actor.Start()
}

// SuperviseSupervisor adds a nested supervisor (creating a hierarchy)
func (s *Supervisor) SuperviseSupervisor(subSupervisor *Supervisor) {
	fmt.Println("Supervisor supervising sub-core...")
	subSupervisor.ctx = s.ctx
	subSupervisor.child = true
	subSupervisor.supervisorMonitor.SetOutboundChannel(s.supervisorMonitor.GetInboundChannel())
	s.subSupervisors[subSupervisor.GetID()] = subSupervisor
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

func (s *Supervisor) handleActorFailure(result *ActorResult) {
	actor := s.findActor(result.ID)

	switch result.Action {
	case ACTOR_RESTART:
		fmt.Println("Restarting actor due to critical error...")
		actor.Stop()
		actor.Start()

	case ACTOR_RETRY:
		fmt.Println("Retrying the failed message...")
		actor.Stop()
		actor.SendMessage(result.Message)
		actor.Start()

	case ACTOR_FAIL:
		fmt.Println("Propagating failure to parent core...")
		// If this supervisor is a child, propagate the failure upwards
		if s.child {
			s.reportErrorToParent(result)
		}
	}
}

func (s *Supervisor) reportErrorToParent(result *ActorResult) {
	s.supervisorMonitor.GetOutboundChannel() <- &SupervisorActorResult{
		Action: SUPERVISOR_FAIL,
		Result: result,
	}
}

func (s *Supervisor) handleSupervisorFailure(result *SupervisorActorResult) {

	switch result.Action {
	case SUPERVISOR_RESTART:
		fmt.Println("Restarting actor due to critical error...")
		// TODO restat supervisor

	case SUPERVISOR_FAIL:
		fmt.Println("Propagating failure to parent core...")
		// TODO RELAY TO NEXT LAYER
	}
}

// findActor finds an actor by name in the supervisor's map
func (s *Supervisor) findActor(id uuid.UUID) Actor {
	if actor, exists := s.actors[id]; exists {
		return actor
	}
	return nil // TODO: Handle if actor is not found, maybe return an error
}
