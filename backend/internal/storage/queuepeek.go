package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/redis/go-redis/v9"
)

// queue stats

type QueuePeekStore struct {
	client *redis.Client
}

func NewQueuePeekStore(client *redis.Client) *QueuePeekStore {
	return &QueuePeekStore{client: client}
}

func (qs *QueuePeekStore) PeekReady(ctx context.Context, count int64) ([]model.Task, error) {
	if count <= 0 {
		count = 20
	}

	results, err := qs.client.ZRangeWithScores(ctx,KeyReady,0,count-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to peek ready queue: %w", err)
	}

	tasks := make([]model.Task, 0, len(results))
	for _,z := range results {
		var task model.Task
		if err := json.Unmarshal([]byte(z.Member.(string)), &task); err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (qs *QueuePeekStore) PeekDelayed(ctx context.Context, count int64) ([]model.DelayedEntry, error) {
	if count <= 0 {
		count = 20
	}

	results, err := qs.client.ZRangeWithScores(ctx,KeyDelayed,0,count-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to peek delayed queue: %w", err)
	}

	tasks := make([]model.DelayedEntry, 0, len(results))
	for _,z := range results {
		var task model.Task
		if err := json.Unmarshal([]byte(z.Member.(string)), &task); err != nil {
			continue
		}
		tasks = append(tasks,model.DelayedEntry{
			Task : task,
			ExecuteAt: z.Score,
		})
	}
	return tasks, nil
}

func (qs *QueuePeekStore) DelayedSize(ctx context.Context) (int64, error) {
	return qs.client.ZCard(ctx, KeyDelayed).Result()
}

func (qs *QueuePeekStore) DeadLetterSize(ctx context.Context) (int64, error) {
	return qs.client.LLen(ctx,KeyDeadLetter).Result()
}