package worker

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime/metrics"
	"time"

	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/xd-sarthak/dtq/internal/queue"
	"github.com/xd-sarthak/dtq/internal/storage"
)

// core processing logic
type Executor struct {
	delayed *queue.DelayedScheduler
	deadLetter *queue.DeadLetterQueue
	metrics *storage.MetricsStore
	events *storage.EventStore
	workerState *storage.WorkerStateStore
}

func NewExecutor(
	delayed *queue.DelayedScheduler,
	deadLetter *queue.DeadLetterQueue,
	metrics *storage.MetricsStore,
	events *storage.EventStore,
	workerState *storage.WorkerStateStore,
) *Executor {
	return &Executor{
		delayed: delayed,
		deadLetter: deadLetter,
		metrics: metrics,
		events: events,
		workerState: workerState,
	}
}

func (e *Executor) Execute(ctx context.Context, workerID int, task model.Task) {
	log.Printf("Worker %d: Starting task %s (priority: %d, retries: %d)", workerID, task.ID, task.Priority, task.Retries)

	// update worker state to busy
	e.workerState.Set(ctx, model.WorkerState{
		WorkerID: workerID,
		Status:   "busy",
		TaskID:   task.ID,
		StartedAt: time.Now(),
	})

	//emit task started event
	e.emitEvent(ctx, task.ID, "started", workerID, fmt.Sprintf("Task %s started execution", task.ID))

	// simulate work by sleeping for a random duration between 100-500ms
	workDuration := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(workDuration)

	// randomly decide if task succeeds or fails (80% success rate)
	if rand.Float64() < 0.2 {
		//failure
		errMsg := fmt.Sprintf("Worker %d: Task %s failed after %v", workerID, task.ID, workDuration)
		log.Println(errMsg)
		task.Error = errMsg
		e.emitEvent(ctx, task.ID, "failed", workerID, errMsg)
		e.handleFailure(ctx, workerID, task)
	} else {
		//success
		log.Printf("Worker %d: Task %s completed successfully in %v", workerID, task.ID, workDuration)
		e.emitEvent(ctx, task.ID, "completed", workerID, fmt.Sprintf("Task %s completed successfully in %v", task.ID, workDuration))
		if err := e.metrics.Increment(ctx); err != nil {
			log.Printf("Failed to increment metrics: %v", err)
		}
	}

	//set worker state back to idle
	e.workerState.Set(ctx, model.WorkerState{
		WorkerID: workerID,
		Status:   "idle",
	})

}

// handle task failure with retry logic
func (e *Executor) handleFailure(ctx context.Context, workerID int, task model.Task) {
	// retry with exponential backoff if retries are left
	if task.Retries < task.MaxRetries {
		task.Retries++
		task.Status = model.StatusPending
		backoff := time.Duration(math.Pow(2, float64(task.Retries))) * time.Second

		log.Printf("Worker %d: Retrying task %s in %v (retry #%d)", workerID, task.ID, backoff, task.Retries)

		e.emitEvent(ctx, task.ID, "retrying", workerID, fmt.Sprintf("Retrying task %s in %v (retry #%d)", task.ID, backoff, task.Retries))

		if err := e.metrics.IncrementRetries(ctx); err != nil {
			log.Printf("Failed to increment retry metrics: %v", err)
		}

		if err := e.delayed.Schedule(ctx, task, backoff); err != nil {
			log.Printf("Failed to schedule task %s for retry: %v", task.ID, err)
		}

	} else {
		// no retries left, move to dead letter queue
		task.Status = model.StatusFailed
		log.Printf("Worker %d: Task %s has reached max retries and will be moved to dead letter queue", workerID, task.ID)

		e.emitEvent(ctx, task.ID, "dead_lettered", workerID, fmt.Sprintf("Task %s has reached max retries and is moved to dead letter queue", task.ID))
		
		failtask := model.FailedTask{
			Task: task,
			FailedAt: time.Now(),
			Error: task.Error,
		}

		if err := e.deadLetter.Add(ctx, failtask); err != nil {
			log.Printf("Failed to add task %s to dead letter queue: %v", task.ID, err)
		}

		if err := e.metrics.IncrementDeadLetters(ctx); err != nil {
			log.Printf("Failed to increment dead letter metrics: %v", err)
		}
	}
}


// emit event for task lifecycle changes
func (e *Executor) emitEvent(ctx context.Context, taskID, eventType string, workerID int, detail string) {
	event := model.TaskEvent{
		ID:        fmt.Sprintf("%s-%d", taskID, time.Now().UnixNano()),
		TaskID:    taskID,
		Type:      eventType,
		WorkerID:  workerID,
		Detail:    detail,
		Timestamp: time.Now(),
	}
	if err := e.events.PushEvent(ctx, event); err != nil {
		log.Printf("Failed to push event for task %s: %v", taskID, err)
	}
}

