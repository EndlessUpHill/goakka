package core

// SupervisorMonitor manages communication between supervisors
type SupervisorMonitor struct {
	supervisor *Supervisor
	inbound    chan *SupervisorActorResult
	outbound   chan *SupervisorActorResult
}

// NewSupervisorMonitor creates a new monitor for supervisors
func NewSupervisorMonitor(s *Supervisor) *SupervisorMonitor {
	return &SupervisorMonitor{
		supervisor: s,
	}
}

func (m *SupervisorMonitor) GetInboundChannel() chan *SupervisorActorResult {
	if m.inbound == nil {
		m.inbound = make(chan *SupervisorActorResult)
		m.monitor()
	}
	return m.inbound
}

func (m *SupervisorMonitor) GetOutboundChannel() chan *SupervisorActorResult {
	return m.outbound
}

func (m *SupervisorMonitor) SetOutboundChannel(outbound chan *SupervisorActorResult) {
	m.outbound = outbound
}

func (m *SupervisorMonitor) monitor() {
	go func() {
		for result := range m.inbound {
			m.supervisor.handleSupervisorFailure(result)
		}
	}()
}

// ActorMonitor manages communication between actors and their supervisor
type ActorMonitor struct {
	supervisor *Supervisor
	inbound    chan *ActorResult
}

// NewActorMonitor creates a new monitor for actors
func NewActorMonitor(s *Supervisor) *ActorMonitor {
	return &ActorMonitor{
		supervisor: s,
	}
}

func (m *ActorMonitor) GetInboundChannel() chan *ActorResult { // Corrected channel type
	if m.inbound == nil {
		m.inbound = make(chan *ActorResult) // Corrected type
		m.monitor()
	}
	return m.inbound
}

func (m *ActorMonitor) monitor() {
	go func() {
		for result := range m.inbound {
			m.supervisor.handleActorFailure(result)
		}
	}()
}
