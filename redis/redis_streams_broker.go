package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EndlessUpHill/goakka/core"
	"github.com/go-redis/redis/v8"
)

// RedisStreamsBroker is an implementation of the MessageBroker interface using Redis Streams
type RedisStreamsBroker struct {
	client     *redis.Client
	ctx        context.Context
	cancel     context.CancelFunc
	groupName  string
	consumerID string
}

// NewRedisStreamsBroker creates a new Redis Streams broker
func NewRedisStreamsBroker(redisAddr, groupName, consumerID string) *RedisStreamsBroker {
	fmt.Println("Creating new Redis Streams broker...")
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx, cancel := context.WithCancel(context.Background())

	return &RedisStreamsBroker{
		client:     client,
		ctx:        ctx,
		cancel:     cancel,
		groupName:  groupName,
		consumerID: consumerID,
	}
}

// Publish sends a message to a Redis stream
func (b *RedisStreamsBroker) Publish(stream string, msg interface{}) error {
	id, err := b.client.XAdd(b.ctx, &redis.XAddArgs{
		Stream: stream,
		Values: msg,
	}).Result()
	if err != nil {
		log.Printf("Error adding message to stream %s: %v", stream, err)
		return err
	}
	fmt.Printf("Message added to stream %s with ID %s\n", stream, id)
	return nil
}

// Subscribe subscribes an actor to a Redis stream group
func (b *RedisStreamsBroker) Subscribe(stream string, actor core.Actor) error {
	// Create the consumer group if it doesn't exist
	err := b.client.XGroupCreateMkStream(b.ctx, stream, b.groupName, "0").Err()
	if err != nil && err != redis.Nil {
		log.Printf("Error creating group %s on stream %s: %v", b.groupName, stream, err)
		return err
	}

	// Process messages in a separate goroutine
	go func() {
		for {
			select {
			case <-b.ctx.Done():
				fmt.Println("Subscription for stream", stream, "has been cancelled.")
				return
			default:
				// Read messages from the stream using the consumer group
				entries, err := b.client.XReadGroup(b.ctx, &redis.XReadGroupArgs{
					Group:    b.groupName,
					Consumer: b.consumerID,
					Streams:  []string{stream, ">"},
					Count:    1,
					Block:    5 * time.Second, // Block for 5 seconds if no message
				}).Result()

				if err != nil && err != redis.Nil {
					log.Printf("Error reading message from stream %s: %v", stream, err)
					continue
				}

				for _, entry := range entries {
					for _, msg := range entry.Messages {

						// Send the message to the actor
						actor.SendMessage(msg.Values)

						// Acknowledge the message after processing
						err = b.client.XAck(b.ctx, stream, b.groupName, msg.ID).Err()
						if err != nil {
							log.Printf("Error acknowledging message %s in stream %s: %v", msg.ID, stream, err)
						} else {
							fmt.Printf("Message %s acknowledged in stream %s\n", msg.ID, stream)
						}
					}
				}
			}
		}
	}()

	return nil
}

// Close gracefully stops the Redis Streams broker and cancels all subscriptions
func (b *RedisStreamsBroker) Close() {
	// Cancel the context to stop all subscription goroutines
	b.cancel()

	// Close the Redis client connection
	b.client.Close()
}
