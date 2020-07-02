package job

import (
	pb "baseGo/src/imserver/api/logic/grpc"
	"baseGo/src/imserver/internal/job/conf"
	"context"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"

	log "fecho/golog"

	cluster "github.com/bsm/sarama-cluster"
)

// Job is push job.
type Job struct {
	c            *conf.Config
	consumer     *cluster.Consumer
	cometServers map[string]*Comet

	rooms      map[string]*Room
	roomsMutex sync.RWMutex
}

// New new a push job.
func New(c *conf.Config) *Job {
	j := &Job{
		c:        c,
		consumer: newKafkaSub(c.Kafka),
		rooms:    make(map[string]*Room),
	}
	go j.newAddress()
	return j
}

func newKafkaSub(c *conf.Kafka) *cluster.Consumer {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(c.Brokers, "chat-im-job", []string{c.Topic}, config)
	//consumer, err := cluster.NewConsumer(c.Brokers, c.Group, []string{"kafka-12"}, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

// Close close resounces.
func (j *Job) Close() error {
	if j.consumer != nil {
		return j.consumer.Close()
	}
	return nil
}

// Consume messages, watch signals
func (j *Job) Consume() {
	for {
		select {
		case err := <-j.consumer.Errors():
			log.Error("Job", "Consume", "consumer error(%v)", err)
		case n := <-j.consumer.Notifications():
			log.Info("Job", "Consume", "consumer rebalanced(%v)", n)
		case msg, ok := <-j.consumer.Messages():
			if !ok {
				log.Info("Job", "Consume", "consumer Message error")
				return
			}
			j.consumer.MarkOffset(msg, "")
			// process push message
			pushMsg := new(pb.PushMsg)
			if err := proto.Unmarshal(msg.Value, pushMsg); err != nil {
				log.Error("Job", "Consume", "proto.Unmarshal(%v) error(%v)", nil, msg, err)
				continue
			}
			if err := j.push(context.Background(), pushMsg); err != nil {
				log.Error("Job", "Consume", "j.push(%v) error(%v)", nil, pushMsg, err)
			}
			//log.Info( "Job", "Consume", "consume: %s/%d/%d\t%s\t%+v", msg.Topic, msg.Partition, msg.Offset, msg.Key, pushMsg)
		}
	}
}

func (j *Job) newAddress() {
	for true {
		j.cometServers = NewComet(j.c.Comet)
		time.Sleep(time.Second * 3)
	}
}
