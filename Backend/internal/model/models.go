package model

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID         string     `json:"id"`
	Priority   int        `json:"priority"`
	Delay      int        `json:"delay"`       
	MaxRetries int        `json:"max_retries"` // 0 = no retries
	Retries    int        `json:"retries"`
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	Error      string     `json:"error,omitempty"`
}

type TaskEvent struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	Type      string    `json:"type"` // submitted, started, completed, failed, retrying, dead_lettered, promoted
	WorkerID  int       `json:"worker_id"`
	Detail    string    `json:"detail"`
	Timestamp time.Time `json:"timestamp"`
}