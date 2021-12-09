package mysql

import "time"

func WithMaxLifetime(duration time.Duration) Option {
	return func(sql *MySQL) {
		sql.maxLifetime = duration
	}
}

func WithMaxOpenConnections(connections int) Option {
	return func(sql *MySQL) {
		sql.maxOpenCons = connections
	}
}

func WithMaxIdleConnections(connections int) Option {
	return func(sql *MySQL) {
		sql.maxIdleCons = connections
	}
}
