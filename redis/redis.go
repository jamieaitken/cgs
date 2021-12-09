package redis

import (
	"github.com/go-redis/redis/v8"
	instr "github.com/jamieaitken/promred/redis"
)

const (
	defaultMasterName = "mymaster"
	defaultPassword   = ""
	defaultDB         = 0
)

type Redis struct {
	masterName string
	addrs      []string
	password   string
	db         int
	clientFunc ClientFunc
	client     instr.Redis
}

type Option func(*Redis)

type ClientFunc func(r *Redis) *redis.Client

func FailOver(r *Redis) *redis.Client {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    r.masterName,
		SentinelAddrs: r.addrs,
		Password:      r.password,
		DB:            r.db,
	})
}

func NonFailOver(r *Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.addrs[0],
		Password: r.password,
		DB:       r.db,
	})
}

func New(addrs []string, opts ...Option) *Redis {
	r := &Redis{
		addrs:      addrs,
		masterName: defaultMasterName,
		password:   defaultPassword,
		db:         defaultDB,
		clientFunc: FailOver,
	}

	r.add(opts...)

	client := r.clientFunc(r)

	r.client = instr.New(client)

	return r
}

func (r *Redis) add(opts ...Option) {
	for _, opt := range opts {
		opt(r)
	}
}

func (r *Redis) Client() instr.Redis {
	return r.client
}

func (r *Redis) Addrs() []string {
	return r.addrs
}

func (r *Redis) MasterName() string {
	return r.masterName
}

func (r *Redis) Password() string {
	return r.password
}

func (r *Redis) DB() int {
	return r.db
}
