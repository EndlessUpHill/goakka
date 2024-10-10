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
	cancel context.CancelFunc // To cancel the subscription goroutines
}

// NewRedisBroker creates a new Redis broker
func NewRedisBroker(redisAddr string) *RedisBroker {
	fmt.Println("Creating new Redis broker...")
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx, cancel := context.WithCancel(context.Background())

	return &RedisBroker{
		client: client,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Publish sends a message to a Redis Pub/Sub topic
func (b *RedisBroker) Publish(topic string, msg interface{}) error {
	err := b.client.Publish(b.ctx, topic, msg).Err()
	if err != nil {
		log.Printf("Error publishing message to topic %s: %v", topic, err)
	}
	return err
}

// Subscribe subscribes an actor to a Redis Pub/Sub topic
func (b *RedisBroker) Subscribe(topic string, actor core.Actor) error {
	sub := b.client.Subscribe(b.ctx, topic)

	// Process messages in a separate goroutine
	go func() {
		defer sub.Close()

		for {
			select {
			case <-b.ctx.Done():
				// Handle context cancellation (shutdown)
				fmt.Println("Subscription for topic", topic, "has been cancelled.")
				return

			default:
				// Receive messages from Redis Pub/Sub
				msg, err := sub.ReceiveMessage(b.ctx)
				if err != nil {
					log.Printf("Error receiving message from topic %s: %v", topic, err)
					return
				}

				// Send the message to the actor
				actor.SendMessage(msg.Payload)
			}
		}
	}()

	return nil
}

// Close gracefully stops the Redis broker and cancels all subscriptions
func (b *RedisBroker) Close() {
	// Cancel the context to stop all subscription goroutines
	b.cancel()

	// Close the Redis client connection
	b.client.Close()
}
