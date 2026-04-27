package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	KeyReady      = "taskqueue:ready"
	KeyDelayed    = "taskqueue:delayed"
	KeyDeadLetter = "taskqueue:deadletter"
	KeyMetrics    = "taskqueue:metrics"
	KeyEvents     = "taskqueue:events"
	KeyWorkers    = "taskqueue:workers"
)

func NewRedisClient(addr, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully")
	return client
}