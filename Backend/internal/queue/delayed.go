package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/xd-sarthak/dtq/internal/storage"
)

type DelayedScheduler struct {
	client *redis.Client
	queue *PriorityQueue
	events *storage.EventStore
}

// initialize the delayed scheduler
func NewDelayedScheduler(client *redis.Client, queue *PriorityQueue, events *storage.EventStore) *DelayedScheduler {
	return &DelayedScheduler{
		client: client,
		queue: queue,
		events: events,
	}
}


// Schedule a task to be executed after a specified delay
func (ds *DelayedScheduler) Schedule(ctx context.Context, task model.Task, delay time.Duration) error {
	data,err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	executeAt := time.Now().Add(delay).Unix()
	return ds.client.ZAdd(ctx,storage.KeyDelayed,redis.Z{
		Score: float64(executeAt),
		Member: data,
	}).Err()
}

// poll the delayed tasks and move them to the main queue when their execution time arrives
func (ds *DelayedScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Delayed scheduler stopped")
			return
		case <-ticker.C:
			ds.completeDueTasks(ctx)
		}
	}
}

// completeDueTasks checks for tasks that are due for execution and moves them to the main queue
func (ds *DelayedScheduler) completeDueTasks(ctx context.Context) {
	now := float64(time.Now().Unix())
	results, err := ds.client.ZRangeByScoreWithScores(ctx, storage.KeyDelayed, &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%f", now),
		Count: 100,
	}).Result()

	if err != nil {
		log.Printf("Failed to fetch due tasks: %v", err)
		return
	}

	for _,z := range results {
		member := z.Member.(string)

		//remove the task from the delayed set
		removed, err := ds.client.ZRem(ctx, storage.KeyDelayed, member).Result()

		if err != nil || removed == 0 {
			continue // maybe another worker has already processed this task
		}

		var task model.Task
		if err := json.Unmarshal([]byte(member), &task); err != nil {
			log.Printf("Failed to unmarshal task: %v", err)
			continue
		}

		// add the task to the main queue
		if err := ds.queue.Enqueue(ctx, task); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
		} else {
			log.Printf("Moved task %s to main queue", task.ID)

			//emit task scheduled event
			if ds.events != nil {
				event := model.TaskEvent {
					ID: fmt.Sprintf("%s-%d", task.ID, time.Now().UnixNano()),
					TaskID: task.ID,
					Type: "scheduled",
					WorkerID: -1, // scheduler is not a worker
					Detail: fmt.Sprintf("Task %s moved to main queue", task.ID),
					Timestamp: time.Now(),
				}
				if err := ds.events.PushEvent(ctx, event); err != nil {
					log.Printf("Failed to push event: %v", err)
				}
			}
		}
	}


}

