package worker

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xd-sarthak/dtq/internal/queue"
	"github.com/xd-sarthak/dtq/internal/storage"
)
// basically a worker pool manager


type WorkerPool struct {
	queue  *queue.PriorityQueue
	executor *Executor
	workercount int
	pollInterval time.Duration
	activeCount atomic.Int64
	wg sync.WaitGroup
	workerState *storage.WorkerStateStore
}

// init a worker pool
func NewWorkerPool(
	queue *queue.PriorityQueue,
	workerState *storage.WorkerStateStore,
	workercount int,
	pollInterval time.Duration,
	executor *Executor)*WorkerPool {
	return &WorkerPool{
		queue: queue,
		executor: executor,
		workercount: workercount,
		pollInterval: pollInterval,
		workerState: workerState,
	}
}

// start the worker pool
func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workercount; i++ {
		wp.wg.Add(1)
		wp.workerState.Set(ctx, storage.WorkerIdleState(i))
		go wp.worker(ctx, i)
	}
	log.Printf("Worker pool started with %d workers", wp.workercount)
}

// wait for all workers to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	log.Println("All workers have stopped")
}

// current active workers
func (wp *WorkerPool) ActiveWorkers() int64 {
	return wp.activeCount.Load()
}

// worker loop
func (wp *WorkerPool) worker(ctx context.Context, workerID int){
	defer wp.wg.Done()
	log.Printf("Worker %d started", workerID)

	for {
		select {
			case <-ctx.Done():
				log.Printf("Worker %d stopping due to context cancellation", workerID)
				return
			default:
		}

		task, err := wp.queue.Dequeue(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("worker %d has dequeue error: %v",workerID,err)
			time.Sleep(wp.pollInterval)
			continue
		}

		if task == nil {
			//queue is empty, wait and retry
			time.Sleep(wp.pollInterval)
			continue
		}

		wp.activeCount.Add(1)
		wp.executor.Execute(ctx, workerID, task)
		wp.activeCount.Add(-1)
	}
}
	