const BASE_URL = "http://localhost:8080";

// ── Type definitions (mirror Go models) ──────────────────────────────

export interface Task {
  id: string;
  priority: number;
  delay: number;
  max_retries: number;
  retries: number;
  status: "pending" | "running" | "completed" | "failed";
  created_at: string;
  error?: string;
}

export interface TaskEvent {
  id: string;
  task_id: string;
  type:
    | "submitted"
    | "started"
    | "completed"
    | "failed"
    | "retrying"
    | "dead_lettered"
    | "promoted";
  worker_id: number;
  detail: string;
  timestamp: string;
}

export interface WorkerState {
  worker_id: number;
  status: "idle" | "busy" | "offline";
  task_id?: string;
  started_at?: string;
}

export interface FailedTask {
  task: Task;
  failed_at: string;
  error: string;
}

export interface EnhancedMetrics {
  total_processed: number;
  total_failed: number;
  total_retries: number;
  queue_size: number;
  active_workers: number;
  success_rate: number;
  delayed_queue_size: number;
  dead_letter_size: number;
  total_submitted: number;
}

export interface DelayedEntry {
  task: Task;
  execute_at: number;
}

export interface QueueData {
  ready: Task[];
  delayed: DelayedEntry[];
}

// ── Fetch helpers ────────────────────────────────────────────────────

async function fetchJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, init);
  if (!res.ok) throw new Error(`API ${path} returned ${res.status}`);
  return res.json() as Promise<T>;
}

export function getMetrics() {
  return fetchJSON<EnhancedMetrics>("/api/metrics/enhanced");
}

export function getEvents(limit = 80) {
  return fetchJSON<TaskEvent[]>(`/api/events?limit=${limit}`);
}

export function getWorkers() {
  return fetchJSON<WorkerState[]>("/api/workers");
}

export function getQueues() {
  return fetchJSON<QueueData>("/api/queues");
}

export function getFailedTasks(offset = 0, limit = 50) {
  return fetchJSON<FailedTask[]>(
    `/api/tasks/failed?offset=${offset}&limit=${limit}`
  );
}

export function submitTask(payload: {
  id: string;
  priority: number;
  delay: number;
  max_retries: number;
}) {
  return fetchJSON<{ message: string }>("/api/tasks", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
}

export function flushData() {
  return fetchJSON<{ message: string }>("/api/flush", { method: "DELETE" });
}
