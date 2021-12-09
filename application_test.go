package cgs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jamieaitken/cgs"
	"github.com/jamieaitken/cgs/config"
	"github.com/jamieaitken/cgs/mysql"
	"github.com/jamieaitken/cgs/publisher"
	"github.com/jamieaitken/cgs/redis"
	"github.com/jamieaitken/cgs/server"
	"github.com/jamieaitken/cgs/subscriber"
	"github.com/jamieaitken/cgs/testing/opts"
	"go.uber.org/zap"
)

func TestNew_Success(t *testing.T) {
	tests := []struct {
		name               string
		givenOpts          []cgs.Option
		expectedLogger     *zap.Logger
		expectedConfig     *config.Config
		expectedRedis      *redis.Redis
		expectedMysql      *mysql.MySQL
		expectedServer     *server.Server
		expectedPublisher  *publisher.KafkaPublisher
		expectedSubscriber *subscriber.KafkaSubscriber
	}{
		{
			name: "given one of each dependency, expect them to be available in the container",
			givenOpts: []cgs.Option{
				cgs.WithLoggerOpts(zap.Development()),
				cgs.WithRedis(context.Background(), "test", []string{"test"}),
				cgs.WithRouter(),
				cgs.WithConfig(config.WithConfigFile("localsettings.env")),
				cgs.WithMySQL(context.Background(), "test", "root:hunter@(localhost:3306)/mysql?parseTime=true"),
				cgs.WithServer(),
				cgs.WithPublisher(context.Background(), "test", []string{"test"}, "test"),
				cgs.WithSubscriber(context.Background(), "test", []string{"test"}, "test"),
			},
			expectedConfig:     loadConfig(t, "localsettings.env"),
			expectedRedis:      redis.New([]string{"test"}),
			expectedMysql:      loadMySQL(t, "root:hunter@(localhost:3306)/mysql?parseTime=true"),
			expectedServer:     server.New(nil),
			expectedPublisher:  loadPublisher(t, []string{"test"}, "test"),
			expectedSubscriber: loadSubscriber(t, []string{"test"}, "test"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New(test.givenOpts...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(app.Config().File(), test.expectedConfig.File()) {
				t.Fatalf(cmp.Diff(app.Config().File(), test.expectedConfig.File()))
			}

			rc, err := app.Redis("test")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(rc, test.expectedRedis, opts.RedisComparer) {
				t.Fatalf(cmp.Diff(rc, test.expectedConfig, opts.RedisComparer))
			}

			sc, err := app.MySQL("test")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(sc, test.expectedMysql, opts.SQLComparer) {
				t.Fatalf(cmp.Diff(sc, test.expectedMysql, opts.SQLComparer))
			}

			if !cmp.Equal(app.Server(), test.expectedServer, opts.ServerComparer) {
				t.Fatalf(cmp.Diff(app.Server(), test.expectedServer, opts.ServerComparer))
			}

			sub, err := app.Subscriber("test")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(sub, test.expectedSubscriber, opts.SubscriberComparer) {
				t.Fatalf(cmp.Diff(sub, test.expectedSubscriber, opts.SubscriberComparer))
			}

			pub, err := app.Publisher("test")
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			if !cmp.Equal(pub, test.expectedPublisher, opts.PublisherComparer) {
				t.Fatalf(cmp.Diff(pub, test.expectedPublisher, opts.PublisherComparer))
			}
		})
	}
}

func TestRedis_Fail(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "given no redis instantiation, expect error when try to access",
			expectedError: cgs.ErrInvalidRedisClient,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New()
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			_, err = app.Redis("blah")
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestMySQL_Fail(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "given no MySQL instantiation, expect error when try to access",
			expectedError: cgs.ErrInvalidMySQLClient,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New()
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			_, err = app.MySQL("blah")
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestPublisher_Fail(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "given no publisher instantiation, expect error when try to access",
			expectedError: cgs.ErrInvalidPublisher,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New()
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			_, err = app.Publisher("blah")
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestSubscriber_Fail(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "given no subscriber instantiation, expect error when try to access",
			expectedError: cgs.ErrInvalidSubscriber,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New()
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			_, err = app.Subscriber("blah")
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestApplication_Run_Success(t *testing.T) {
	tests := []struct {
		name       string
		givenFuncs []func(ctx context.Context) error
	}{
		{
			name: "given an error, expect it to be raised",
			givenFuncs: []func(ctx context.Context) error{
				func(ctx context.Context) error {
					return nil
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New()
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			err = app.Run(context.Background(), test.givenFuncs...)
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}
}

func TestApplication_Run_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenFuncs    []func(ctx context.Context) error
		expectedError error
	}{
		{
			name: "given an error, expect it to be raised",
			givenFuncs: []func(ctx context.Context) error{
				func(ctx context.Context) error {
					return errors.New("func errored")
				},
			},
			expectedError: cgs.ErrRunFailed,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := cgs.New()
			if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}

			err = app.Run(context.Background(), test.givenFuncs...)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatalf(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func loadConfig(t *testing.T, file string) *config.Config {
	c, err := config.New(config.WithConfigFile(file))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	return c
}

func loadMySQL(t *testing.T, addr string) *mysql.MySQL {
	s, err := mysql.New(addr)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	return s
}

func loadPublisher(t *testing.T, addr []string, topic string, ops ...publisher.Option) *publisher.KafkaPublisher {
	p, err := publisher.New(addr, topic, ops...)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	return p
}

func loadSubscriber(t *testing.T, addr []string, topic string, ops ...subscriber.Option) *subscriber.KafkaSubscriber {
	p, err := subscriber.New(addr, topic, ops...)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	return p
}
