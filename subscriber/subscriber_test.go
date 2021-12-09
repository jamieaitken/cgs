package subscriber_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/segmentio/kafka-go"

	"github.com/google/go-cmp/cmp"
	"github.com/jamieaitken/cgs/subscriber"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                string
		givenAddrs          []string
		givenTopic          string
		givenOpts           []subscriber.Option
		expectedAddr        []string
		expectedTopic       string
		expectedMaxAttempts int
	}{
		{
			name:                "given custom max attempts, expect default values for everything else",
			givenAddrs:          []string{"10.00.00.1"},
			givenTopic:          "test",
			givenOpts:           []subscriber.Option{subscriber.WithMaxAttempts(40)},
			expectedAddr:        []string{"10.00.00.1"},
			expectedTopic:       "test",
			expectedMaxAttempts: 40,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := subscriber.New(test.givenAddrs, test.givenTopic, test.givenOpts...)
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

			if !cmp.Equal(actual.Topic(), test.expectedTopic) {
				t.Fatalf(cmp.Diff(actual.Topic(), test.expectedTopic))
			}
		})
	}
}

func TestKafkaSubscriber_Ping_Success(t *testing.T) {
	tests := []struct {
		name      string
		addr      []string
		topic     string
		givenOpts []subscriber.Option
	}{
		{
			name:  "given success, expect zero errors",
			addr:  []string{"10.0.0.1"},
			topic: "test",
			givenOpts: []subscriber.Option{subscriber.WithHealthChecker(&mockHealthChecker{
				GivenResponse: &kafka.MetadataResponse{},
			})},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := subscriber.New(test.addr, test.topic, test.givenOpts...)
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

func TestKafkaSubscriber_Ping_Fail(t *testing.T) {
	tests := []struct {
		name          string
		addr          []string
		topic         string
		givenOpts     []subscriber.Option
		expectedError error
	}{
		{
			name:  "given health checker error, expect error to be raised",
			addr:  []string{"10.0.0.1"},
			topic: "test",
			givenOpts: []subscriber.Option{subscriber.WithHealthChecker(&mockHealthChecker{
				GivenError: errors.New("failed to create request"),
			})},
			expectedError: subscriber.ErrFailedToContactBroker,
		},
		{
			name:  "given topic error, expect error to be raised",
			addr:  []string{"10.0.0.1"},
			topic: "test",
			givenOpts: []subscriber.Option{subscriber.WithHealthChecker(&mockHealthChecker{
				GivenResponse: &kafka.MetadataResponse{
					Topics: []kafka.Topic{
						{},
						{
							Error: errors.New("fail"),
						},
					},
				},
			})},
			expectedError: subscriber.ErrFailedToReadTopic,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s, err := subscriber.New(test.addr, test.topic, test.givenOpts...)
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
