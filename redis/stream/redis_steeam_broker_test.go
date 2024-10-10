package stream_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/EndlessUpHill/goakka/core"
	coreRedis "github.com/EndlessUpHill/goakka/redis"
	"github.com/ory/dockertest"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var redisClient *redis.Client
var redisStreamBroker *coreRedis.RedisStreamsBroker
var redisBroker *coreRedis.RedisBroker

// TestMain is called before and after the test suite
func TestMain(m *testing.M) {
	if os.Getenv("CI") == "true" {
		// Running in CI (GitHub Actions), connect to Redis service
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}

		redisClient = redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})

		if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
			log.Fatalf("Could not connect to Redis: %s", err)
		}

		redisBroker = coreRedis.NewRedisBroker(redisAddr)
		redisStreamBroker = coreRedis.NewRedisStreamsBroker(redisAddr, "test-group", "test-consumer")
	} else {
		// Running locally, use dockertest to start Redis
		var pool *dockertest.Pool
		var resource *dockertest.Resource
		var err error

		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}

		resource, err = pool.Run("redis", "6.2", nil)
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}

		// Connect to Redis
		if err := pool.Retry(func() error {
			redisClient = redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
			})
			_, err := redisClient.Ping(context.Background()).Result()
			return err
		}); err != nil {
			log.Fatalf("Could not connect to Redis: %s", err)
		}


		redisBroker = coreRedis.NewRedisBroker(fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")))
		redisStreamBroker = coreRedis.NewRedisStreamsBroker(fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")), "test-group", "test-consumer")

		// Cleanup Redis container after tests
		defer func() {
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("Could not purge resource: %s", err)
			}
		}()
	}

	// Run the tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

func TestRedisStreamPubSub(t *testing.T) {
    // Create a test actor
    var wg sync.WaitGroup
    wg.Add(1)
    var msg map[string]interface{}

    actor := core.NewBasicActor("test-actor", func(res *core.ActorResult) *core.ActorResult {
        // Store the received message and ensure it's cast correctly
        if receivedMsg, ok := res.Message.(map[string]interface{}); ok {
            msg = receivedMsg
        } else {
            t.Fatalf("Received message is not in expected format: %v", res.Message)
        }

        wg.Done() // Mark the wait group as done only after processing
        return &core.ActorResult{}
    })
	actor.Start()

    // Subscribe the actor to the Redis stream
    err := redisStreamBroker.Subscribe("test-stream", actor)
    assert.NoError(t, err)

    // Publish a test message to the Redis stream
    testMessage := map[string]interface{}{
        "field1": "value1",
        "field2": "value2",
    }
    err = redisStreamBroker.Publish("test-stream", testMessage)
    assert.NoError(t, err)

    // Wait for the actor to receive the message
    wg.Wait()

    // Assert that the received message matches the published message
    assert.Equal(t, testMessage["field1"], msg["field1"])
    assert.Equal(t, testMessage["field2"], msg["field2"])

	 // Delay the cancellation of the Redis context to allow acknowledgment
    time.Sleep(1 * time.Second)
	
    // Clean up
    redisStreamBroker.Close()
}