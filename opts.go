package cgs

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/jamieaitken/cgs/config"
	"github.com/jamieaitken/cgs/mysql"
	"github.com/jamieaitken/cgs/publisher"
	"github.com/jamieaitken/cgs/redis"
	"github.com/jamieaitken/cgs/router"
	"github.com/jamieaitken/cgs/server"
	"github.com/jamieaitken/cgs/subscriber"
	"github.com/spf13/viper"
)

const (
	registeredMsg = "registered %s-%s"
)

func WithLoggerOpts(opts ...zap.Option) Option {
	return func(application *Application) {
		l, err := zap.NewProduction(opts...)
		if err != nil {
			application.logger.Error("failed to apply new logger", zap.Error(err))
			return
		}

		application.logger = l
	}
}

func WithRedis(ctx context.Context, name string, addrs []string, opts ...redis.Option) Option {
	return func(application *Application) {
		application.mu.Lock()
		defer application.mu.Unlock()

		r := redis.New(addrs, opts...)
		application.redis[name] = r

		application.health.AddReadinessCheck(fmt.Sprintf("%s-redis", name), func() error {
			err := r.Client().Ping(ctx, "application").Err()
			if err != nil {
				application.logger.Error(fmt.Sprintf("%s-redis failed healthcheck", name), zap.Error(err))
			}

			return err
		})
		application.logger.Info(fmt.Sprintf(registeredMsg, name, "redis"))
	}
}

func WithMySQL(ctx context.Context, name, addr string, opts ...mysql.Option) Option {
	return func(application *Application) {
		application.mu.Lock()
		defer application.mu.Unlock()

		m, err := mysql.New(addr, opts...)
		if err != nil {
			application.logger.Error(fmt.Sprintf("failed to create sql client for %s", name), zap.Error(err))
			return
		}

		application.mysql[name] = m
		application.health.AddReadinessCheck(fmt.Sprintf("%s-mysql", name), func() error {
			err = m.Client().PingContext(ctx)
			if err != nil {
				application.logger.Error(fmt.Sprintf("%s-mysql failed healthcheck", name), zap.Error(err))
			}

			return err
		})
		application.logger.Info(fmt.Sprintf(registeredMsg, name, "mysql"))
	}
}

func WithPublisher(ctx context.Context, name string, addrs []string, topic string, opts ...publisher.Option) Option {
	return func(application *Application) {
		application.mu.Lock()
		defer application.mu.Unlock()

		p, err := publisher.New(addrs, topic, opts...)
		if err != nil {
			application.logger.Error(fmt.Sprintf("failed to create kafka publisher for %s", name), zap.Error(err))
			return
		}

		application.publishers[name] = p
		application.health.AddReadinessCheck(fmt.Sprintf("%s-publisher", name), func() error {
			err = p.Ping(ctx)
			if err != nil {
				application.logger.Error(fmt.Sprintf("%s-publisher failed healthcheck", name), zap.Error(err))
			}

			return err
		})
		application.logger.Info(fmt.Sprintf(registeredMsg, name, "publisher"))
	}
}

func WithSubscriber(ctx context.Context, name string, addrs []string, topic string, opts ...subscriber.Option) Option {
	return func(application *Application) {
		application.mu.Lock()
		defer application.mu.Unlock()

		p, err := subscriber.New(addrs, topic, opts...)
		if err != nil {
			application.logger.Error(fmt.Sprintf("failed to create kafka subscriber for %s", name), zap.Error(err))
			return
		}

		application.subscribers[name] = p
		application.health.AddReadinessCheck(fmt.Sprintf("%s-subscriber", name), func() error {
			err = p.Ping(ctx)
			if err != nil {
				application.logger.Error(fmt.Sprintf("%s-subscriber failed healthcheck", name), zap.Error(err))
			}

			return err
		})
		application.logger.Info(fmt.Sprintf(registeredMsg, name, "subscriber"))
	}
}

func WithRouter(opts ...router.Option) Option {
	return func(application *Application) {
		if application.router == nil {
			application.router = router.New(application.health, opts...)

			application.logger.Info("registered new router")

			return
		}

		application.router.Add(opts...)

		application.logger.Info("altered registered router")
	}
}

func WithConfig(opts ...config.Option) Option {
	return func(application *Application) {
		cfg, err := config.New(opts...)
		if err != nil {
			application.logger.Error("failed to register new config", zap.Error(err))

			return
		}

		if cfg.File() != "" {
			application.logger.Info(fmt.Sprintf("using config file %s", viper.ConfigFileUsed()))
		}

		application.config = cfg

		application.logger.Info("registered new config")
	}
}

func WithServer(opts ...server.Option) Option {
	return func(application *Application) {
		application.server = server.New(application.router.Mux(), opts...)

		application.logger.Info("registered new server")
	}
}
