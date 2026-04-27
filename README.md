# Distributed Task Queue

A high-performance, educational distributed task queue built with Go, Redis, and Next.js. This system is designed to provide visibility into distributed task processing, showcasing a robust backend architecture paired with a real-time reactive dashboard.

## 🚀 Features

- **Robust Go Backend:** Implements an at-least-once delivery queue using Go and Redis with features like:
  - Priority queuing
  - Delayed task scheduling
  - Retries & Dead Letter Queues (DLQ)
  - Worker pool management
- **Real-time Next.js Dashboard:** A beautifully designed frontend to visualize the whole pipeline.
  - Live task flow pipeline
  - Real-time worker status (idle/busy/offline)
  - Activity logs tracking exactly what stage tasks are in
  - High-density metrics and success rates
- **Docker Ready:** Effortless local setup using Docker Compose.

## 🏗 Architecture

- **Backend:** Go (Golang) + go-redis. Exposes a REST API on `localhost:8080`.
- **Database/Broker:** Redis 7 (stores pending queues, delayed sets, worker states, and metrics).
- **Frontend:** Next.js 16 + React 19 + TailwindCSS v4. Runs on `localhost:3000`.

## ⚙️ Quick Start

### 1. Start the Backend and Redis

Use Docker Compose to spin up the Go backend and Redis instance:

```bash
cd dist-task-queue
docker-compose up -d
```
*The api will be available at `http://localhost:8080`.*

### 2. Start the Frontend Dashboard

Ensure you have Node.js 20+ installed.

```bash
cd frontend/my-app
npm install  # or pnpm install / yarn install
npm run dev
```
*Open your browser to `http://localhost:3000` to access the dashboard.*

## 📖 API Usage Example

Submit a task directly to the backend bypassing the UI:

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"id": "test-task", "priority": 5, "delay": 2, "max_retries": 3}'
```

## 🛠 Project Structure

```
dist-task-queue/
├── backend/            # Go backend service
│   ├── cmd/server/     # Go application entrypoint
│   └── internal/       # Worker pool, redis storage, and API logic
├── frontend/           # Next.js React Dashboard
│   └── my-app/         # Frontend app logic & components
└── docker-compose.yml  # Docker orchestration for Backend + Redis
```

## 📄 License

MIT License
