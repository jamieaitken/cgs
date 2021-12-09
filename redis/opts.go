package redis

func WithMasterName(masterName string) Option {
	return func(redis *Redis) {
		redis.masterName = masterName
	}
}

func WithPassword(password string) Option {
	return func(redis *Redis) {
		redis.password = password
	}
}

func WithDB(db int) Option {
	return func(redis *Redis) {
		redis.db = db
	}
}

func WithClientType(clientFunc ClientFunc) Option {
	return func(redis *Redis) {
		redis.clientFunc = clientFunc
	}
}
