package subscriber

import (
	"context"
	"errors"
	"fmt"
	instr "github.com/jamieaitken/promred/kafka"
	"github.com/segmentio/kafka-go"
)

var (
	ErrFailedToAssertKafkaClient = errors.New("failed to assert health checker to kafka client")
	ErrFailedToReadTopic         = errors.New("failed to read kafka topic")
	ErrFailedToContactBroker     = errors.New("failed to contact broker")
)

const (
	defaultMaxAttempts = 5
)

type KafkaSubscriber struct {
	addrs         []string
	topic         string
	maxAttempts   int
	healthChecker HealthChecker
	client        instr.Reader
}

type Option func(*KafkaSubscriber)

func New(addrs []string, topic string, opts ...Option) (*KafkaSubscriber, error) {
	k := &KafkaSubscriber{
		addrs:         addrs,
		topic:         topic,
		maxAttempts:   defaultMaxAttempts,
		healthChecker: &kafka.Client{},
	}

	k.add(opts...)

	r := &kafka.ReaderConfig{
		Brokers:     k.addrs,
		Topic:       k.topic,
		MaxAttempts: k.maxAttempts,
	}

	k.client = instr.NewReader(kafka.NewReader(*r))

	return k, nil
}

func (k *KafkaSubscriber) add(opts ...Option) {
	for _, opt := range opts {
		opt(k)
	}
}

func (k *KafkaSubscriber) Client() instr.Reader {
	return k.client
}

func (k *KafkaSubscriber) Addrs() []string {
	return k.addrs
}

func (k *KafkaSubscriber) Topic() string {
	return k.topic
}

func (k *KafkaSubscriber) MaxAttempts() int {
	return k.maxAttempts
}

func (k *KafkaSubscriber) Ping(ctx context.Context) error {
	res, err := k.healthChecker.Metadata(ctx, &kafka.MetadataRequest{
		Addr:   kafka.TCP(k.addrs...),
		Topics: []string{k.topic},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", err, ErrFailedToContactBroker)
	}

	for _, topic := range res.Topics {
		if topic.Error == nil {
			continue
		}

		return fmt.Errorf("%s: %w", err, ErrFailedToReadTopic)
	}

	return nil
}
