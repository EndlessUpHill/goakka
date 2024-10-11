package main

import (
	"fmt"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/EndlessUpHill/goakka/redis"
)

var RedisRegistryInstance *core.ActorRegistry
var RedisBrokerInstance *redis.RedisBroker

func init() {
	RedisRegistryInstance = core.NewActorRegistry()
	RedisBrokerInstance = redis.NewRedisBroker("localhost:6379")
}

type ExampleRedis struct {
	core.BasicActor
}

func (a *ExampleRedis) ReceiveMessage(msg interface{}) *core.ActorResult {
	fmt.Printf("ExampleCoew Actor %s received message: %v\n", a.GetID(), msg)
	if RedisRegistryInstance != nil {
		actor, found := RedisRegistryInstance.GetActor("actor4")
		if found {
			actor.SendMessage("Hello from ExampleCore")
			RedisBrokerInstance.Publish("example", "Hello from ExampleCore. I broadcasted this message.")
		}
	}
	return &core.ActorResult{}
}

func NewExampleRedis(id string) *ExampleRedis {

	return &ExampleRedis{
		BasicActor: *core.NewBasicActor(id),
	}

}
