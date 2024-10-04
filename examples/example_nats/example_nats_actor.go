package main

import (
	"fmt"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/EndlessUpHill/goakka/nats"
)

var NatsRegistryInstance *core.ActorRegistry 
var NatsBrokerInstance *nats.NatsBroker

func init() {
	NatsRegistryInstance = core.NewActorRegistry()
	NatsBrokerInstance = nats.NewNatsBroker("nats://localhost:4222")
}

type ExampleNats struct {
	core.BasicActor
}
func (a *ExampleNats) ReceiveMessage(msg interface{}) *core.ActorResult {
	fmt.Printf("ExampleCoew Actor %s received message: %v\n", a.GetID(), msg)
	if NatsRegistryInstance != nil {
		actor, found := NatsRegistryInstance.GetActor("actor4")
		if found {
			actor.SendMessage("Hello from ExampleCore")
			NatsBrokerInstance.Publish("example", "Hello from ExampleCore. I broadcasted this message.")
		}
	}
	return &core.ActorResult{}
}

func NewExampleNats(id string) *ExampleNats {

	return &ExampleNats{
		BasicActor: *core.NewBasicActor(id, nil),
	}

}
