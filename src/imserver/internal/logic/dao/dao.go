package dao

import (
	"baseGo/src/imserver/internal/logic/conf"
	"baseGo/src/imserver/internal/logic/model"

	"github.com/gomodule/redigo/redis"

	kafka "gopkg.in/Shopify/sarama.v1"
)

// Dao dao.
type Dao struct {
	c           *conf.Config
	kafkaPub    kafka.SyncProducer
	redis       *redis.Pool
	redisExpire int32
}

// New new a dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:           c,
		kafkaPub:    newKafkaPub(c.Kafka),
		redis:       model.GetRedis(),
		redisExpire: 600,
	}
	return d
}

func newKafkaPub(c *conf.Kafka) kafka.SyncProducer {
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

// Close close the resource.
func (d *Dao) Close() error {
	return d.redis.Close()
}
