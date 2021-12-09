package publisher

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type HealthChecker interface {
	Metadata(ctx context.Context, req *kafka.MetadataRequest) (*kafka.MetadataResponse, error)
}

func WithMaxAttempts(attempts int) Option {
	return func(publisher *KafkaPublisher) {
		publisher.maxAttempts = attempts
	}
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(publisher *KafkaPublisher) {
		checker, ok := publisher.healthChecker.(*kafka.Client)
		if !ok {
			return
		}

		publisher.writeTimeout = timeout

		checker.Timeout = timeout

		publisher.healthChecker = checker
	}
}

func WithRequiredAck(ra kafka.RequiredAcks) Option {
	return func(publisher *KafkaPublisher) {
		publisher.requiredAck = ra
	}
}

func WithHealthChecker(healthChecker HealthChecker) Option {
	return func(publisher *KafkaPublisher) {
		publisher.healthChecker = healthChecker
	}
}
