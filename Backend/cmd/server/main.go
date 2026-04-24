package main

import (
	"fmt"
	"sync"
	"time"
	"container/heap"
	"crypto/rand"
)

// status type
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

//task model
type Task struct {
	ID 		string
	Payload string
	Priority int
	Status   Status
	Attempts int
	MaxRetries int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// make a new id
func newID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//here we created a new task with the given payload, priority and max retries. The task is initialized with a unique ID, pending status, and timestamps for creation and update.
func NewTask(payload string, priority, maxRetries int) Task {
	now := time.Now()
	return Task{
		ID:		 newID(),
		Payload:  payload,
		Priority: priority,
		Status:   StatusPending,
		Attempts: 0,
		MaxRetries: maxRetries,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// task queue -> using heap

// heap in Go
type taskHeap []Task

func (h taskHeap) Len() int           { return len(h) }
func (h taskHeap) Less(i, j int) bool { return h[i].Priority > h[j].Priority } // max-heap
func (h taskHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *taskHeap) Push(x any) {
	*h = append(*h, x.(Task))
}

func (h *taskHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// the queue

// mutex is used to ensure that only one goroutine can access the queue at a time
type Queue struct {
	mu sync.Mutex
	data taskHeap
}

func NewQueue() *Queue {
	q := &Queue{}
	heap.Init(&q.data)
	return q
}

func (q *Queue) Enqueue(task Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(&q.data, task)
}

func (q *Queue) Dequeue() (Task, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.data) == 0 {
		return Task{}, false
	}

	t := heap.Pop(&q.data).(Task)
	return t, true
}

func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.data)
}


// ----------main code---------------

func main() {
	q := NewQueue()
 
	tasks := []Task{
		NewTask("send welcome email",       1, 3),
		NewTask("process payment",          5, 5), // highest priority
		NewTask("generate monthly report",  2, 2),
		NewTask("resize uploaded image",    3, 3),
		NewTask("sync data to warehouse",   4, 4),
	}
 
	fmt.Println("=== Enqueueing tasks ===")
	for _, t := range tasks {
		q.Enqueue(t)
		fmt.Printf("  [+] %-30s  priority=%d  id=%s\n", t.Payload, t.Priority, t.ID[:8])
	}
 
	fmt.Printf("\nQueue size after enqueue: %d\n", q.Size())
 
	fmt.Println("\n=== Dequeueing tasks (expect high→low priority) ===")
	order := 1
	for {
		t, ok := q.Dequeue()
		if !ok {
			break
		}
		fmt.Printf("  %d. %-30s  priority=%d  status=%-10s  id=%s\n",
			order, t.Payload, t.Priority, t.Status, t.ID[:8])
		order++
	}
 
	fmt.Printf("\nQueue size after dequeue: %d\n", q.Size())
}