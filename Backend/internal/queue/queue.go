package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/xd-sarthak/dtq/internal/storage"
)

type PriorityQueue struct {
	client *redis.Client
}

func NewPriorityQueue(client *redis.Client) *PriorityQueue {
	return &PriorityQueue{client: client}
}

// enqueue a task to the ready queue
func (pq *PriorityQueue) Enqueue(ctx context.Context, task model.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	return pq.client.ZAdd(ctx,storage.KeyReady, redis.Z{
		Score: float64(-task.Priority), // redis sorted set uses score for ordering, we use negative priority to get higher priority tasks first
		Member: string(data),
	}).Err()
}

// dequeue pops the highest priority task from the ready queue
func (pq *PriorityQueue) Dequeue(ctx context.Context) (*model.Task, error) {
	
	result, err := pq.client.ZPopMin(ctx, storage.KeyReady, 1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to pop task: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no tasks available")
	}

	var task model.Task
	if err := json.Unmarshal([]byte(result[0].Member.(string)), &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// size returns the number of tasks in the ready queue
func (pq *PriorityQueue) Size(ctx context.Context) (int64, error) {
	return pq.client.ZCard(ctx, storage.KeyReady).Result()
}