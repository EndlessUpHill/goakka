package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/EndlessUpHill/goakka/core"
	coreRedis "github.com/EndlessUpHill/goakka/redis"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisStreamPubSub(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Make sure Redis is running
	})
	defer client.Close()

	// Create the Redis Stream Pub/Sub system
	streamName := "test-stream"
	groupName := "test-group"
	consumerName := "test-consumer"
	redisStreamPubSub := coreRedis.NewRedisStreamPubSub("localhost:6379", streamName, groupName, consumerName)

	// Publish a message
	testMessage := "Hello, Redis Stream!"

	// Create an actor
	actor := core.NewBasicActor("test-actor", func(msg interface{}) *core.ActorResult {
		assert.Equal(t, "Hello, Redis Stream!", msg)
		return &core.ActorResult{}
	})
	actor.Start()

	// Subscribe the actor to the stream
	redisStreamPubSub.Subscribe(actor)
	redisStreamPubSub.Publish(testMessage)

	// Wait a bit for the message to be processed
	time.Sleep(100 * time.Millisecond)

	// Read the message from the stream and validate
	results, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{streamName, ">"},
		Count:    1,
		Block:    time.Second,
	}).Result()

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, testMessage, results[0].Messages[0].Values["message"])

	// Clean up by deleting the stream
	client.Del(ctx, streamName)
}
