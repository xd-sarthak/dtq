package api

import "net/http"

func NewRouter(handler *APIHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/tasks", handler.SubmitTask)
	mux.HandleFunc("GET /api/metrics", handler.GetMetrics)
	mux.HandleFunc("GET /api/tasks/failed", handler.GetFailedTasks)
	mux.HandleFunc("GET /api/health", handler.HealthCheck)

	mux.HandleFunc("GET /api/events", handler.RecentActivity)
	mux.HandleFunc("GET /api/workers", handler.RuntimeState)
	mux.HandleFunc("GET /api/queues", handler.GetQueues)
	mux.HandleFunc("GET /api/metrics/enhanced", handler.GetEnhancedMetrics)
	mux.HandleFunc("DELETE /api/flush", handler.FlushData)

	var apihandler http.Handler = mux
	apihandler = corsMiddleware(loggingMiddleware(recoveryMiddleware(apihandler)))

	return apihandler
}