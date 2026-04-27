package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/xd-sarthak/dtq/internal/api"
	"github.com/xd-sarthak/dtq/internal/queue"
	"github.com/xd-sarthak/dtq/internal/storage"
	"github.com/xd-sarthak/dtq/internal/worker"
	"github.com/xd-sarthak/dtq/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Redis client
	redisClient := storage.NewRedisClient(cfg.RedisAddr,cfg.RedisPassword)
	defer redisClient.Close()

	// Initialize storage components
	queuePeekStore := storage.NewQueuePeekStore(redisClient)
	workerStateStore := storage.NewWorkerStateStore(redisClient)
	deadLetterQueue := storage.NewDeadLetterQueue(redisClient)
	metricsStore := storage.NewMetricsStore(redisClient)
	eventStore := storage.NewEventStore(redisClient)

	// queue
	priorityQueue := queue.NewPriorityQueue(redisClient)
	delayedScheduler := queue.NewDelayedScheduler(redisClient,priorityQueue,eventStore)

	//workers
	executor := worker.NewExecutor(delayedScheduler,deadLetterQueue,metricsStore,eventStore,workerStateStore)
	pool := worker.NewWorkerPool(priorityQueue,workerStateStore,cfg.WorkerCount,time.Duration(cfg.PollInterval)*time.Second,executor)

	// api
	apiHandler := api.NewAPIHandler(priorityQueue,delayedScheduler,deadLetterQueue,metricsStore,workerStateStore,pool,redisClient,eventStore,queuePeekStore)
	router := api.NewRouter(apiHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// graceful shutdown setup
	ctx,stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go delayedScheduler.Start(ctx)

	go pool.Start(ctx)

    go func(){
		log.Printf("Server is running on port %d", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	} ()

	// wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received, shutting down gracefully...")

	// shutdown server with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	pool.Wait()
	log.Println("Server exited gracefully")
}

