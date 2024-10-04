package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)
 

type ActorResult struct {
	Error      error       
	Action     int     
	Message    interface{} 
	name       string
	id		 uuid.UUID
}

type Actor interface {
	Start()
	Stop()
	ReceiveMessage(msg interface{}) *ActorResult
	SendMessage(msg interface{})
	GetID() uuid.UUID
	GetName() string
	SetWaitGroup(wg *sync.WaitGroup)
	SetContext(ctx context.Context)
	SetFailureChannel(chan *ActorResult)

}

type BasicActor struct {
	id             uuid.UUID
	name		   string
	mailbox     chan interface{}
	stop        chan struct{}
	wg          *sync.WaitGroup
	ctx         context.Context
	receiveFunc func(result *ActorResult) *ActorResult
	failureChannel chan *ActorResult
}

func NewBasicActor(name string, recieveFunc func(result *ActorResult) *ActorResult) *BasicActor {
	return &BasicActor{
		id: uuid.New(),
		name: name,
		mailbox:     make(chan interface{}, 100),
		stop:        make(chan struct{}),
		receiveFunc: recieveFunc,
	}
}

func (a *BasicActor) GetID() uuid.UUID {
	return a.id
}

func (a *BasicActor) GetName() string {
	return a.name
}

func (a *BasicActor) SetWaitGroup(wg *sync.WaitGroup) {
	a.wg = wg
}

func (a *BasicActor) SetContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *BasicActor) SetFailureChannel(failure chan *ActorResult) {
	a.failureChannel = failure
}

func (a *BasicActor) Start() {
	fmt.Printf("Starting actor %s...\n", a.id)
	if a.wg != nil {
		a.wg.Add(1)
	}
	go func() {
		defer func() {
			fmt.Printf("Actor %s finished.\n", a.id)
			if a.wg != nil {
				a.wg.Done()
			}
		}()
		for {
			select {
			case msg := <-a.mailbox:
				var result *ActorResult
				if a.receiveFunc != nil {
					actor := ActorResult{
						Message: msg,
						name: a.name,
						id: a.id,
					}
					result = a.receiveFunc(&actor)
				} else {
					result = a.ReceiveMessage(msg)
				}

				if result.Error != nil {
					fmt.Printf("Actor %s encountered a failure: %v\n", a.GetID(), result.Error)
					a.failureChannel <- result
				}
			case <-a.stop:
				fmt.Printf("Stopping actor %s due to stop signal.\n", a.id)
				return
			case <-a.ctx.Done():
				fmt.Printf("Stopping actor %s due to context cancellation.\n", a.id)
				return
			}
		}
	}()
}

func (a *BasicActor) Stop() {
	fmt.Printf("Stopping actor %s...\n", a.id)
	select {
	case <-a.stop:
		// Already closed
	default:
		close(a.stop)
	}
}

func (a *BasicActor) ReceiveMessage(msg interface{}) *ActorResult {
	fmt.Printf("!!!Actor received message: %v\n", msg)
	return &ActorResult{}
}

func (a *BasicActor) SendMessage(msg interface{}) {
	select {
	case a.mailbox <- msg:
	default:
		fmt.Printf("Actor %s mailbox full, dropping message: %v\n", a.id, msg)
	}
}
