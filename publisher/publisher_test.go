package publisher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jamieaitken/cgs/publisher"
	"github.com/jamieaitken/cgs/testing/opts"
	"github.com/segmentio/kafka-go"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                  string
		givenAddrs            []string
		givenTopic            string
		givenOpts             []publisher.Option
		expectedAddr          []string
		expectedTopic         string
		expectedMaxAttempts   int
		expectedWriteTimeout  time.Duration
		expectedRequiredAck   kafka.RequiredAcks
		expectedHealthChecker *kafka.Client
	}{
		{
			name:                 "given custom max attempts, expect default values for everything else",
			givenAddrs:           []string{"10.00.00.1"},
			givenTopic:           "test",
			givenOpts:            []publisher.Option{publisher.WithMaxAttempts(40)},
			expectedAddr:         []string{"10.00.00.1"},
			expectedTopic:        "test",
			expectedMaxAttempts:  40,
			expectedWriteTimeout: time.Second * 20,
			expectedRequiredAck:  kafka.RequireAll,
			expectedHealthChecker: &kafka.Client{
				Addr:    kafka.TCP([]string{"10.00.00.1"}...),
				Timeout: time.Second * 20,
			},
		},
		{
			name:                 "given custom write Timeout, expect default values for everything else",
			givenAddrs:           []string{"10.00.00.1"},
			givenTopic:           "test",
			givenOpts:            []publisher.Option{publisher.WithWriteTimeout(time.Second * 400)},
			expectedAddr:         []string{"10.00.00.1"},
			expectedTopic:        "test",
			expectedMaxAttempts:  5,
			expectedWriteTimeout: time.Second * 400,
			expectedRequiredAck:  kafka.RequireAll,
			expectedHealthChecker: &kafka.Client{
				Addr:    kafka.TCP([]string{"10.00.00.1"}...),
				Timeout: time.Second * 400,
			},
		},
		{
			name:                 "given custom required ack, expect default values for everything else",
			givenAddrs:           []string{"10.00.00.1"},
			givenTopic:           "test",
			givenOpts:            []publisher.Option{publisher.WithRequiredAck(kafka.RequireNone)},
			expectedAddr:         []string{"10.00.00.1"},
			expectedTopic:        "test",
			expectedMaxAttempts:  5,
			expectedWriteTimeout: time.Second * 20,
			expectedRequiredAck:  kafka.RequireNone,
			expectedHealthChecker: &kafka.Client{
				Addr:    kafka.TCP([]string{"10.00.00.1"}...),
				Timeout: time.Second * 20,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := publisher.New(test.givenAddrs, test.givenTopic, test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(actual.Addrs(), test.expectedAddr) {
				t.Fatalf(cmp.Diff(actual.Addrs(), test.expectedAddr))
			}

			if !cmp.Equal(actual.MaxAttempts(), test.expectedMaxAttempts) {
				t.Fatalf(cmp.Diff(actual.MaxAttempts(), test.expectedMaxAttempts))
			}

			if !cmp.Equal(actual.Addrs(), test.expectedAddr) {
				t.Fatalf(cmp.Diff(actual.Addrs(), test.expectedAddr))
			}

			if !cmp.Equal(actual.WriteTimeout(), test.expectedWriteTimeout) {
				t.Fatalf(cmp.Diff(actual.WriteTimeout(), test.expectedWriteTimeout))
			}

			if !cmp.Equal(actual.RequiredAck(), test.expectedRequiredAck) {
				t.Fatalf(cmp.Diff(actual.RequiredAck(), test.expectedRequiredAck))
			}

			if !cmp.Equal(actual.HealthChecker(), test.expectedHealthChecker, opts.KafkaClientComparer) {
				t.Fatalf(cmp.Diff(actual.HealthChecker(), test.expectedHealthChecker, opts.KafkaClientComparer))
			}
		})
	}
}

func TestKafkaPublisher_Ping_Success(t *testing.T) {
	tests := []struct {
		name      string
		addr      []string
		topic     string
		givenOpts []publisher.Option
	}{
		{
			name:  "given success, expect zero errors",
			addr:  []string{"10.0.0.1"},
			topic: "test",
			givenOpts: []publisher.Option{publisher.WithHealthChecker(&mockHealthChecker{
				GivenResponse: &kafka.MetadataResponse{},
			})},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := publisher.New(test.addr, test.topic, test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			err = s.Ping(context.Background())
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}
}

func TestKafkaPublisher_Ping_Fail(t *testing.T) {
	tests := []struct {
		name          string
		addr          []string
		topic         string
		givenOpts     []publisher.Option
		expectedError error
	}{
		{
			name:  "given health checker error, expect error to be raised",
			addr:  []string{"10.0.0.1"},
			topic: "test",
			givenOpts: []publisher.Option{publisher.WithHealthChecker(&mockHealthChecker{
				GivenError: errors.New("failed to create request"),
			})},
			expectedError: publisher.ErrFailedToContactBroker,
		},
		{
			name:  "given topic error, expect error to be raised",
			addr:  []string{"10.0.0.1"},
			topic: "test",
			givenOpts: []publisher.Option{publisher.WithHealthChecker(&mockHealthChecker{
				GivenResponse: &kafka.MetadataResponse{
					Topics: []kafka.Topic{
						{},
						{
							Error: errors.New("fail"),
						},
					},
				},
			})},
			expectedError: publisher.ErrFailedToReadTopic,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := publisher.New(test.addr, test.topic, test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			err = s.Ping(context.Background())
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockHealthChecker struct {
	GivenResponse *kafka.MetadataResponse
	GivenError    error
}

func (m *mockHealthChecker) Metadata(_ context.Context, _ *kafka.MetadataRequest) (*kafka.MetadataResponse, error) {
	return m.GivenResponse, m.GivenError
}
