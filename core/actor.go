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
}

type BasicActor struct {
	mailbox chan interface{}
	stop    chan struct{}
	wg      *sync.WaitGroup
	ctx     context.Context
}

func NewBasicActor() *BasicActor {
	return &BasicActor{
		mailbox: make(chan interface{}, 100),
		stop:    make(chan struct{}),
	}
}

func (a *BasicActor) setWaitGroup(wg *sync.WaitGroup) {
	a.wg = wg
}

func (a *BasicActor) setContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *BasicActor) Start() {
	fmt.Println("Starting actor...")
	if a.wg != nil {
		a.wg.Add(1)
	}
	go func() {
		defer func() {
			fmt.Println("Actor Finished.")
			a.wg.Done()
		}()

		for {
			select {
			case msg := <-a.mailbox:
				a.Receive(msg)
			case <-a.stop:
				fmt.Println("Stopping actor... due to context cancellation")
				return
			}
		}
	}()
}

func (a *BasicActor) Stop() {
	fmt.Println("Stopping actor...")
	// a.stop <- struct{}{}
	close(a.stop)
}

func (a *BasicActor) Receive(msg interface{}) {
	fmt.Printf("Actor received message: %v\n", msg)
}

func (a *BasicActor) SendMessage(msg interface{}) {
	fmt.Printf("Actor sending message: %v\n", msg)
	a.mailbox <- msg
}
