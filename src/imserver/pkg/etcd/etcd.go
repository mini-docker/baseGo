package etcd

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/fecho/registry"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type Master struct {
	Path   string
	Nodes  *sync.Map
	Client *clientv3.Client
}

//node is a client
type Node struct {
	Key  string
	Meta ServiceMeta
}

//the detail of service
type ServiceMeta struct {
	IP  string
	Ext interface{}
}

type Service struct {
	Name    string
	Meta    ServiceMeta
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
}

//func NewMaster(endpoints []string, watchPath string) (*Master, error) {
// Monitor Nodes
func OnWatch(endpoints []string, watchPath string) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	})

	if err != nil {
		golog.Error("etcd", "OnWatch", "", err)
		return nil, err
	}

	master := &Master{
		Path:   watchPath,
		Nodes:  new(sync.Map),
		Client: cli,
	}

	go master.WatchNodes()
	return master, err
}

func (m *Master) AddNode(key string, info *ServiceMeta) {

	node := &Node{
		Key:  key,
		Meta: *info,
	}

	m.Nodes.Store(node.Key, node)
}

func GetServiceMeta(ev *clientv3.Event) *ServiceMeta {
	info := &ServiceMeta{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		golog.Error("etcd", "GetServiceMeta", "", err)
	}
	return info
}

func (m *Master) WatchNodes() {

	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())

	for wresp := range rch {

		for _, ev := range wresp.Events {

			switch ev.Type {
			case clientv3.EventTypePut:
				info := GetServiceMeta(ev)
				m.AddNode(string(ev.Kv.Key), info)

			case clientv3.EventTypeDelete:
				m.Nodes.Delete(string(ev.Kv.Key))
			}
		}
	}
}

// Register Service
func Register(name string, info ServiceMeta, endpoints []string, dialTimeout int) (*Service, error) {
	config := clientv3.Config{}

	if dialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	} else {
		config.DialTimeout = time.Duration(dialTimeout) * time.Second
	}

	if len(endpoints) > 0 {
		config.Endpoints = endpoints
	}
	cli, err := clientv3.New(config)
	if err != nil {
		golog.Error("etcd", "Register", "", err)
		return nil, err
	}
	return &Service{
		Name:   name,
		Meta:   info,
		stop:   make(chan error),
		client: cli,
	}, err
}

func (s *Service) Start() error {
	ch, err := s.keepAlive()
	if err != nil {
		golog.Error("etcd", "Start", "", err)
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				golog.Info("etcd", "Start", "", ka)
				s.revoke()
				return nil
			} else {
				golog.Info("etcd", "Start", "", ka)
			}
		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	info := &s.Meta

	key := "services/" + s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		golog.Error("etcd", "keepAlive", "", err)
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		golog.Error("etcd", "keepAlive", "", err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		golog.Error("etcd", "revoke", "", err)
	}

	//log.Printf("servide:%s stop\n", s.Name)
	return err
}

func InitEtcd(opts ...registry.Option) error {
	config := clientv3.Config{}

	var options registry.Options

	for _, o := range opts {
		o(&options)
	}

	if options.Timeout == 0 {
		options.Timeout = 5 * time.Second
	}

	if options.Secure || options.TLSConfig != nil {
		tlsConfig := options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		config.TLS = tlsConfig
	}

	var cAddrs []string

	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	cli, err := clientv3.New(config)
	if err != nil {
		return err
	}

	registry.SetEtcdClient(cli)
	return nil
}
