package opts

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jamieaitken/cgs/mysql"
	"github.com/jamieaitken/cgs/publisher"
	"github.com/jamieaitken/cgs/redis"
	"github.com/jamieaitken/cgs/server"
	"github.com/jamieaitken/cgs/subscriber"
	"github.com/segmentio/kafka-go"
)

var KafkaClientComparer = cmp.Comparer(func(x, y *kafka.Client) bool {
	return x.Addr.String() == y.Addr.String() &&
		x.Timeout.String() == y.Timeout.String()
})

var RedisComparer = cmp.Comparer(func(x, y redis.Redis) bool {
	return x.DB() == y.DB() && cmp.Equal(x.Addrs(), y.Addrs()) &&
		x.Password() == y.Password() && x.MasterName() == y.MasterName()
})

var SQLComparer = cmp.Comparer(func(x, y mysql.MySQL) bool {
	return x.MaxIdleCons() == y.MaxIdleCons() && x.MaxOpenCons() == y.MaxOpenCons() &&
		x.MaxLifetime() == y.MaxLifetime() && x.Addr() == y.Addr()
})

var ServerComparer = cmp.Comparer(func(x, y server.Server) bool {
	return x.ReadTimeout() == y.ReadTimeout() && x.WriteTimeout() == y.WriteTimeout() &&
		x.Addr() == y.Addr()
})

var SubscriberComparer = cmp.Comparer(func(x, y subscriber.KafkaSubscriber) bool {
	return x.Topic() == y.Topic() && cmp.Equal(x.Addrs(), y.Addrs()) &&
		x.MaxAttempts() == y.MaxAttempts()
})

var PublisherComparer = cmp.Comparer(func(x, y publisher.KafkaPublisher) bool {
	return x.Topic() == y.Topic() && cmp.Equal(x.Addrs(), y.Addrs()) &&
		x.MaxAttempts() == y.MaxAttempts() && x.RequiredAck() == y.RequiredAck()
})
