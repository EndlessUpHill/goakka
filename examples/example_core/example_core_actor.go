package main

import (
	"fmt"

	"github.com/EndlessUpHill/goakka/core"
)

var RegistryInstance *core.ActorRegistry
var BrokerInstance *core.InMemoryBroker

func init() {
	RegistryInstance = core.NewActorRegistry()
	BrokerInstance = core.NewInMemoryBroker()
}

type ExampleCore struct {
	core.BasicActor
}

func (a *ExampleCore) ReceiveMessage(msg interface{}) {
	fmt.Printf("ExampleCore Actor %s received message: %v\n", a.GetID(), msg)
	if RegistryInstance != nil {
		actor, found := RegistryInstance.GetActor("actor1")
		if found {
			actor.SendMessage("Hello from ExampleCore")
			BrokerInstance.Publish("example", "Hello from ExampleCore. I broadcasted this message.")
		}
	}
}

func NewExampleCore(id string) *ExampleCore {

	return &ExampleCore{
		BasicActor: *core.NewBasicActor(id, nil),
	}

}
