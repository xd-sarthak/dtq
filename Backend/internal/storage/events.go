package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/redis/go-redis/v9"
)

type EventStore struct {
	client *redis.Client
}

func NewEventStore(client *redis.Client) *EventStore {
	return &EventStore{client: client}
}

// push a task event to the Redis stream
func (es *EventStore) PushEvent(ctx context.Context, event model.TaskEvent) error {
	data ,err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	pipe := es.client.Pipeline()
	pipe.LPush(ctx, KeyEvents, string(data))
	pipe.LTrim(ctx, KeyEvents, 0, 999)
	_, err = pipe.Exec(ctx)
	return err
}

// list returns the most recent task events, newest first
func (es *EventStore) List(ctx context.Context, count int64) ([]model.TaskEvent, error){
	if count <= 0 {
		count = 50
	}

	results, err := es.client.LRange(ctx, KeyEvents, 0, count-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events: %w", err)
	}

	events := make([]model.TaskEvent, 0, len(results))
	for _, res := range results {
		var event model.TaskEvent
		if err := json.Unmarshal([]byte(res), &event); err != nil {
			continue
		}
		events = append(events, event)
	}
	return events, nil
}