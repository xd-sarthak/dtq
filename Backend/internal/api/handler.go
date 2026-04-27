package api

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/xd-sarthak/dtq/internal/queue"
	"github.com/xd-sarthak/dtq/internal/storage"
	"github.com/xd-sarthak/dtq/internal/worker"

)

type APIHandler struct {
	queue 	*queue.PriorityQueue
	delayed *queue.DelayedScheduler
	deadLetter *storage.DeadLetterQueue
	metrics *storage.MetricsStore
	workerState *storage.WorkerStateStore
	pool *worker.WorkerPool
	redis *redis.Client
	events *storage.EventStore
	queuePeek *storage.QueuePeekStore
}

func NewAPIHandler(
	queue *queue.PriorityQueue,
	delayed *queue.DelayedScheduler,
	deadLetter *storage.DeadLetterQueue,
	metrics *storage.MetricsStore,
	workerState *storage.WorkerStateStore,
	pool *worker.WorkerPool,
	redis *redis.Client,
	events *storage.EventStore,
	queuePeek *storage.QueuePeekStore,
) *APIHandler {
	return &APIHandler{
		queue: queue,
		delayed: delayed,
		deadLetter: deadLetter,
		metrics: metrics,
		workerState: workerState,
		pool: pool,
		redis: redis,
		events: events,
		queuePeek: queuePeek,
	}
}

// incoming payload for task submission
type submitTaskRequest struct {
	ID string `json:"id"`
	Priority int `json:"priority"`
	Delay int `json:"delay"`
	MaxRetries int `json:"max_retries"`
}

// SubmitTask handles incoming task submission requests, validates the payload, and enqueues the task for processing.
func (h *APIHandler) SubmitTask(w http.ResponseWriter, r *http.Request) {
	var req submitTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}

	if req.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "task ID is required"})
		return
	}

	task := model.Task{
		ID: req.ID,
		Priority: req.Priority,
		Delay: req.Delay,
		MaxRetries: req.MaxRetries,
		Status: model.StatusPending,
		CreatedAt: time.Now(),
	}

	ctx := r.Context()

	// emit task submitted event
	event := model.TaskEvent{
	ID: fmt.Sprintf("%s-%d", task.ID, time.Now().UnixNano()),
	TaskID: task.ID,
	Type: "submitted",
	WorkerID: -1, // not assigned to any worker yet
	Detail: fmt.Sprintf("Task %s submitted with priority %d and delay %d seconds", task.ID, task.Priority, task.Delay),
	Timestamp: time.Now(),
	}

	h.events.PushEvent(ctx, event)
	h.metrics.IncrementSubmitted(ctx)

	if req.Delay > 0 {
		delay := time.Duration(req.Delay) * time.Second
		if err := h.delayed.Schedule(ctx, task, delay); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to schedule delayed task"})
			return
		}
	} else {
		if err := h.queue.Enqueue(ctx, task); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to enqueue task"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "task submitted successfully"})
}

// get metrics
func (h *APIHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	queueSize,_ := h.queue.Size(ctx)
	activeWorkers := h.pool.ActiveWorkers()

	m, err := h.metrics.GetMetrics(ctx, queueSize, activeWorkers)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch metrics"})
		return
	}
	writeJSON(w,http.StatusOK,m)
}

// get failed tasks
func (h *APIHandler) GetFailedTasks(w http.ResponseWriter, r *http.Request) {
	offset,_ := strconv.ParseInt(r.URL.Query().Get("offset"),10,64)
	limit,_ := strconv.ParseInt(r.URL.Query().Get("limit"),10,64)
	if limit <= 0 {
		limit = 20
	}

	tasks, err := h.deadLetter.List(r.Context(), offset, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch failed tasks"})
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

// HealthCheck responds with the system's status by pinging the underlying Redis datastore.
func (h *APIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.redis.Ping(r.Context()).Err(); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}




// writeJSON is a helper to consistently format and dispatch JSON responses.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
