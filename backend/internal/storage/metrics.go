package storage

import (
	"context"
	"strconv"
	"fmt"

	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/redis/go-redis/v9"
)
// rw counts of pipeline throughput
type MetricsStore struct {
	client *redis.Client
}

func NewMetricsStore(client *redis.Client) *MetricsStore {
	return &MetricsStore{client: client}
}

// increment the count of processed tasks
func (ms *MetricsStore) IncrementProcessed(ctx context.Context) error {
	return ms.client.HIncrBy(ctx, KeyMetrics, "processed", 1).Err()
}

// increment the count of failed tasks
func (ms *MetricsStore) IncrementFailed(ctx context.Context) error {
	return ms.client.HIncrBy(ctx, KeyMetrics, "failed", 1).Err()
}

// retries count
func (ms *MetricsStore) IncrementRetries(ctx context.Context) error {
	return ms.client.HIncrBy(ctx, KeyMetrics, "retries", 1).Err()
}

// submit count
func (ms *MetricsStore) IncrementSubmitted(ctx context.Context) error {
	return ms.client.HIncrBy(ctx, KeyMetrics, "submitted", 1).Err()
}

// get all metrics
func (ms *MetricsStore) GetMetrics(ctx context.Context, queueSize, activeWorkers int64) (*model.Metrics, error) {

	vals,err := ms.client.HGetAll(ctx, KeyMetrics).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}

	metrics := &model.Metrics{
		TotalProcessed: parseI64(vals["processed"]),
		TotalFailed:    parseI64(vals["failed"]),
		TotalRetries:   parseI64(vals["retries"]),
		QueueSize:     queueSize,
		ActiveWorkers: activeWorkers,
	}

	return metrics, nil
}

// augment metrics with queue size and active worker count
func (ms *MetricsStore) GetAugmentedMetrics(ctx context.Context, queueSize, activeWorkers,delayedSize,deadLetterSize int64) (*model.EnhancedMetrics, error) {
	vals, err := ms.client.HGetAll(ctx, KeyMetrics).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics: %w", err)
	}

	processed := parseI64(vals["processed"])
	failed := parseI64(vals["failed"])
	submitted := parseI64(vals["submitted"])

	var successRate float64
	total := processed + failed
	if total > 0 {
		successRate = float64(processed) / float64(total) * 100
	}

	return &model.EnhancedMetrics{
		Metrics: model.Metrics{
			TotalProcessed: processed,
			TotalFailed:    failed,
			TotalRetries:   parseI64(vals["retries"]),
			QueueSize:      queueSize,
			ActiveWorkers:  activeWorkers,
		},
		SuccessRate:      successRate,
		DelayedQueueSize: delayedSize,
		DeadLetterSize:   deadLetterSize,
		TotalSubmitted:   submitted,
	}, nil
}
func parseI64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}


