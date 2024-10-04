package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/EndlessUpHill/goakka/core"
)

type RedisStreamPubSub struct {
	client     *redis.Client
	streamName string
	groupName  string
	consumer   string
	ctx        context.Context
}

// NewRedisStreamPubSub creates a new Redis stream-based pub/sub system
func NewRedisStreamPubSub(addr, streamName, groupName, consumer string) *RedisStreamPubSub {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx := context.Background()

	// Create the stream and consumer group if they don't exist
	_, err := client.XGroupCreateMkStream(ctx, streamName, groupName, "0").Result()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	return &RedisStreamPubSub{
		client:     client,
		streamName: streamName,
		groupName:  groupName,
		consumer:   consumer,
		ctx:        ctx,
	}
}

// Publish a message to the Redis stream
func (r *RedisStreamPubSub) Publish(msg string) {
	err := r.client.XAdd(r.ctx, &redis.XAddArgs{
		Stream: r.streamName,
		Values: map[string]interface{}{"message": msg},
	}).Err()
	if err != nil {
		log.Printf("Error adding message to stream: %v", err)
	}
}

// Subscribe an actor to the Redis stream via a consumer group
func (r *RedisStreamPubSub) Subscribe(actor *core.BasicActor) {
	go func() {
		for {
			// Read messages from the consumer group
			result, err := r.client.XReadGroup(r.ctx, &redis.XReadGroupArgs{
				Group:    r.groupName,
				Consumer: r.consumer,
				Streams:  []string{r.streamName, ">"},
				Count:    1,
				Block:    time.Second * 5,
			}).Result()

			if err != nil && err != redis.Nil {
				log.Printf("Error reading from stream: %v", err)
				continue
			}

			for _, stream := range result {
				for _, message := range stream.Messages {
					msg := message.Values["message"].(string)
					actor.SendMessage(msg)
					// Acknowledge the message after it's been processed
					_, err := r.client.XAck(r.ctx, r.streamName, r.groupName, message.ID).Result()
					if err != nil {
						log.Printf("Error acknowledging message: %v", err)
					}
				}
			}
		}
	}()
}


