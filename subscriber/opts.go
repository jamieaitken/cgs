package subscriber

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type HealthChecker interface {
	Metadata(ctx context.Context, req *kafka.MetadataRequest) (*kafka.MetadataResponse, error)
}

func WithMaxAttempts(attempts int) Option {
	return func(subscriber *KafkaSubscriber) {
		subscriber.maxAttempts = attempts
	}
}

func WithHealthChecker(healthChecker HealthChecker) Option {
	return func(subscriber *KafkaSubscriber) {
		subscriber.healthChecker = healthChecker
	}
}
