package cgs

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/heptiolabs/healthcheck"
	"github.com/jamieaitken/cgs/config"
	"github.com/jamieaitken/cgs/mysql"
	"github.com/jamieaitken/cgs/publisher"
	"github.com/jamieaitken/cgs/redis"
	"github.com/jamieaitken/cgs/router"
	"github.com/jamieaitken/cgs/server"
	"github.com/jamieaitken/cgs/subscriber"
)

var (
	ErrInvalidRedisClient = errors.New("redis client not found for given key")
	ErrInvalidMySQLClient = errors.New("mysql client not found for given key")
	ErrInvalidPublisher   = errors.New("kafka publisher not found for given key")
	ErrInvalidSubscriber  = errors.New("kafka subscriber not found for given key")
	ErrRunFailed          = errors.New("run errored")
)

type Application struct {
	logger      *zap.Logger
	config      *config.Config
	redis       map[string]*redis.Redis
	mysql       map[string]*mysql.MySQL
	publishers  map[string]*publisher.KafkaPublisher
	subscribers map[string]*subscriber.KafkaSubscriber
	server      *server.Server
	router      *router.Router
	health      healthcheck.Handler
	mu          sync.Mutex
}

type Option func(*Application)

func New(opts ...Option) (*Application, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	app := &Application{
		logger:      l,
		health:      healthcheck.NewHandler(),
		redis:       make(map[string]*redis.Redis),
		mysql:       make(map[string]*mysql.MySQL),
		publishers:  make(map[string]*publisher.KafkaPublisher),
		subscribers: make(map[string]*subscriber.KafkaSubscriber),
	}

	app.Add(opts...)

	return app, nil
}

func (a *Application) Add(opts ...Option) {
	for _, opt := range opts {
		opt(a)
	}
}

func (a *Application) Server() *server.Server {
	return a.server
}

func (a *Application) Router() *router.Router {
	return a.router
}

func (a *Application) Logger() *zap.Logger {
	return a.logger
}

func (a *Application) Config() *config.Config {
	return a.config
}

func (a *Application) Redis(key string) (*redis.Redis, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	c, ok := a.redis[key]
	if !ok {
		a.logger.Error(fmt.Sprintf("failed to get redis client for %s", key), zap.Error(ErrInvalidRedisClient))

		return nil, ErrInvalidRedisClient
	}

	return c, nil
}

func (a *Application) MySQL(key string) (*mysql.MySQL, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	c, ok := a.mysql[key]
	if !ok {
		a.logger.Error(fmt.Sprintf("failed to get mysql client for %s", key), zap.Error(ErrInvalidMySQLClient))

		return nil, ErrInvalidMySQLClient
	}

	return c, nil
}

func (a *Application) Publisher(key string) (*publisher.KafkaPublisher, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	c, ok := a.publishers[key]
	if !ok {
		a.logger.Error(fmt.Sprintf("failed to get kafka publisher for %s", key), zap.Error(ErrInvalidPublisher))

		return nil, ErrInvalidPublisher
	}

	return c, nil
}

func (a *Application) Subscriber(key string) (*subscriber.KafkaSubscriber, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	c, ok := a.subscribers[key]
	if !ok {
		a.logger.Error(fmt.Sprintf("failed to get kafka subscriber for %s", key), zap.Error(ErrInvalidSubscriber))

		return nil, ErrInvalidSubscriber
	}

	return c, nil
}

func (a *Application) Run(ctx context.Context, fns ...func(ctx context.Context) error) error {
	for _, fn := range fns {
		err := fn(ctx)
		if err != nil {
			a.logger.Error("run error", zap.Error(err))

			return fmt.Errorf("%s: %w", err, ErrRunFailed)
		}
	}

	return nil
}
