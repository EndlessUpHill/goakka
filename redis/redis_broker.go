package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/go-redis/redis/v8"
)

// RedisBroker is an implementation of the MessageBroker interface using Redis Pub/Sub
type RedisBroker struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisBroker creates a new Redis broker
func NewRedisBroker(redisAddr string) *RedisBroker {
	fmt.Println("Creating new Redis broker...")
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx := context.Background()

	return &RedisBroker{
		client: client,
		ctx:    ctx,
	}
}

// Publish sends a message to a Redis Pub/Sub topic
func (b *RedisBroker) Publish(topic string, msg interface{}) error {
	return b.client.Publish(b.ctx, topic, msg).Err()
}

// Subscribe subscribes an actor to a Redis Pub/Sub topic
func (b *RedisBroker) Subscribe(topic string, actor core.Actor) error {
	sub := b.client.Subscribe(b.ctx, topic)

	// Process messages in a separate goroutine
	go func() {
		for {
			msg, err := sub.ReceiveMessage(b.ctx)
			if err != nil {
				log.Printf("error receiving message: %v", err)
				return
			}

			actor.SendMessage(msg.Payload)
		}
	}()

	return nil
}
