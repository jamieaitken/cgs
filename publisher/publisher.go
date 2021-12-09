package publisher

import (
	"context"
	"errors"
	"fmt"
	"time"

	instr "github.com/jamieaitken/promred/kafka"
	"github.com/segmentio/kafka-go"
)

var (
	ErrFailedToReadTopic     = errors.New("failed to read kafka topic")
	ErrFailedToContactBroker = errors.New("failed to contact broker")
)

const (
	defaultMaxAttempts  = 5
	defaultWriteTimeout = time.Second * 20
)

type KafkaPublisher struct {
	addrs         []string
	topic         string
	maxAttempts   int
	writeTimeout  time.Duration
	requiredAck   kafka.RequiredAcks
	publisher     instr.Writer
	healthChecker HealthChecker
}

type Option func(*KafkaPublisher)

func New(addrs []string, topic string, opts ...Option) (*KafkaPublisher, error) {
	k := &KafkaPublisher{
		addrs:        addrs,
		topic:        topic,
		maxAttempts:  defaultMaxAttempts,
		writeTimeout: defaultWriteTimeout,
		requiredAck:  kafka.RequireAll,
		healthChecker: &kafka.Client{
			Addr:    kafka.TCP(addrs...),
			Timeout: defaultWriteTimeout,
		},
	}

	k.add(opts...)

	w := &kafka.Writer{
		Addr:         kafka.TCP(k.addrs...),
		Topic:        k.topic,
		MaxAttempts:  k.maxAttempts,
		WriteTimeout: k.writeTimeout,
		RequiredAcks: k.requiredAck,
	}

	k.publisher = instr.NewWriter(w)

	return k, nil
}

func (k *KafkaPublisher) add(opts ...Option) {
	for _, opt := range opts {
		opt(k)
	}
}

func (k *KafkaPublisher) Client() instr.Writer {
	return k.publisher
}

func (k *KafkaPublisher) Addrs() []string {
	return k.addrs
}

func (k *KafkaPublisher) Topic() string {
	return k.topic
}

func (k *KafkaPublisher) MaxAttempts() int {
	return k.maxAttempts
}

func (k *KafkaPublisher) WriteTimeout() time.Duration {
	return k.writeTimeout
}

func (k *KafkaPublisher) RequiredAck() kafka.RequiredAcks {
	return k.requiredAck
}

func (k *KafkaPublisher) HealthChecker() HealthChecker {
	return k.healthChecker
}

func (k *KafkaPublisher) Ping(ctx context.Context) error {
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
