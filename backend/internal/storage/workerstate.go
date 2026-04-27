package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/redis/go-redis/v9"
)



type WorkerStateStore struct {
	client *redis.Client
}

func WorkerIdleState(workerID int) model.WorkerState {
	return model.WorkerState{
		WorkerID: workerID,
		Status:   "idle",
	}
}

// set updates the state of a worker in Redis
func (ws *WorkerStateStore) Set(ctx context.Context, state model.WorkerState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal worker state: %w", err)
	}
	return ws.client.HSet(ctx, KeyWorkers, strconv.Itoa(state.WorkerID), string(data)).Err()
}

// get retrieves the state of a worker from Redis
func (ws *WorkerStateStore) GetAll(ctx context.Context) ([]model.WorkerState, error){
	results, err := ws.client.HGetAll(ctx, KeyWorkers).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch worker states: %w", err)
	}

	states := make([]model.WorkerState, 0, len(results))
	for _, res := range results {
		var state model.WorkerState
		if err := json.Unmarshal([]byte(res), &state); err != nil {
			continue
		}
		states = append(states, state)
	}
	return states, nil
}