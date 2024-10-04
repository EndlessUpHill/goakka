package core

import (
	"context"
	"fmt"
	"sync"
)

type Actor interface {
	Start()
	Stop()
	ReceiveMessage(msg interface{})
	SendMessage(msg interface{})
	GetID() string
	SetWaitGroup(wg *sync.WaitGroup)
	SetContext(ctx context.Context)
}

type RecieveFunc func(msg interface{})
type BasicActor struct {
	id          string
	mailbox     chan interface{}
	stop        chan struct{}
	wg          *sync.WaitGroup
	ctx         context.Context
	receiveFunc RecieveFunc
}

func NewBasicActor(id string, recieveFunc RecieveFunc) *BasicActor {
	return &BasicActor{
		id:          id,
		mailbox:     make(chan interface{}, 100),
		stop:        make(chan struct{}),
		receiveFunc: recieveFunc,
	}
}

func (a *BasicActor) GetID() string {
	return a.id
}

func (a *BasicActor) SetWaitGroup(wg *sync.WaitGroup) {
	a.wg = wg
}

func (a *BasicActor) SetContext(ctx context.Context) {
	a.ctx = ctx
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
				if a.receiveFunc != nil {
					a.receiveFunc(msg)
				} else {
					a.ReceiveMessage(msg)
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

func (a *BasicActor) ReceiveMessage(msg interface{}) {
	fmt.Printf("!!!Actor received message: %v\n", msg)
}

func (a *BasicActor) SendMessage(msg interface{}) {
	select {
	case a.mailbox <- msg:
	default:
		fmt.Printf("Actor %s mailbox full, dropping message: %v\n", a.id, msg)
	}
}
