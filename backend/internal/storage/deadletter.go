package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/redis/go-redis/v9"
)

type DeadLetterQueue struct {
	client *redis.Client
}

func NewDeadLetterQueue(client *redis.Client) *DeadLetterQueue {
	return &DeadLetterQueue{client: client}
}

// push a failed task to the dead letter queue
func (dlq *DeadLetterQueue) Add(ctx context.Context, failedTask model.FailedTask) error {
	data, err := json.Marshal(failedTask)
	if err != nil {
		return fmt.Errorf("failed to marshal failed task: %w", err)
	}

	if err := dlq.client.LPush(ctx, KeyDeadLetter, data).Err(); err != nil {
		return fmt.Errorf("failed to push task to dead letter queue: %w", err)
	}

	return nil
}

// list returns the most recent failed tasks in the dead letter queue
func (dlq *DeadLetterQueue) List(ctx context.Context, offset, count int64) ([]model.FailedTask, error) {
	results, err := dlq.client.LRange(ctx, KeyDeadLetter, offset, offset+count-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dead letter tasks: %w", err)
	}

	failedTasks := make([]model.FailedTask, 0, len(results))
	for _, res := range results {
		var task model.FailedTask
		if err := json.Unmarshal([]byte(res), &task); err != nil {
			continue
		}
		failedTasks = append(failedTasks, task)
	}
	return failedTasks, nil
}