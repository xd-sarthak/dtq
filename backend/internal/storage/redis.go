package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	KeyReady      = "taskqueue:ready"
	KeyDelayed    = "taskqueue:delayed"
	KeyDeadLetter = "taskqueue:deadletter"
	KeyMetrics    = "taskqueue:metrics"
	KeyEvents     = "taskqueue:events"
	KeyWorkers    = "taskqueue:workers"

	redisMaxRetries    = 10
	redisRetryInterval = 2 * time.Second
)

func NewRedisClient(addr, password string) *redis.Client {
	log.Printf("[redis] Initializing client | addr=%s network=tcp4", addr)

	client := redis.NewClient(&redis.Options{
		Network:  "tcp4",
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	var lastErr error

	for attempt := 1; attempt <= redisMaxRetries; attempt++ {
		lastErr = client.Ping(ctx).Err()
		if lastErr == nil {
			log.Printf("[redis] Connected successfully | addr=%s attempt=%d", addr, attempt)
			return client
		}
		log.Printf("[redis] Connection attempt %d/%d failed | addr=%s err=%v | retrying in %v",
			attempt, redisMaxRetries, addr, lastErr, redisRetryInterval)
		time.Sleep(redisRetryInterval)
	}

	log.Fatalf("[redis] Could not connect after %d attempts | addr=%s last_err=%v", redisMaxRetries, addr, lastErr)
	return nil // unreachable, but satisfies the compiler
}