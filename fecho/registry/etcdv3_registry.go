// Package etcdv3 provides an etcd version 3 registry
package registry

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"strings"
	"sync"
	"time"

	hash "github.com/mitchellh/hashstructure"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

var (
	prefix       = "/services"
	etcdClient   *clientv3.Client
	etcdClientRW sync.RWMutex
)

func SetEtcdClient(c *clientv3.Client) {
	etcdClientRW.Lock()
	defer etcdClientRW.Unlock()
	etcdClient = c
}

func getEtcdClient() *clientv3.Client {
	etcdClientRW.RLock()
	defer etcdClientRW.RUnlock()
	return etcdClient
}

func GetEtcdClient() *clientv3.Client {
	return getEtcdClient()
}

type etcdv3Registry struct {
	client  *clientv3.Client
	options Options
	sync.Mutex
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
}

func configure(e *etcdv3Registry, opts ...Option) error {
	//config := clientv3.Config{}

	for _, o := range opts {
		o(&e.options)
	}

	if e.options.Timeout == 0 {
		e.options.Timeout = 5 * time.Second
	}
	//
	//if e.options.Secure || e.options.TLSConfig != nil {
	//	tlsConfig := e.options.TLSConfig
	//	if tlsConfig == nil {
	//		tlsConfig = &tls.Config{
	//			InsecureSkipVerify: true,
	//		}
	//	}
	//
	//	config.TLS = tlsConfig
	//}
	//
	//var cAddrs []string
	//
	//for _, addr := range e.options.Addrs {
	//	if len(addr) == 0 {
	//		continue
	//	}
	//	cAddrs = append(cAddrs, addr)
	//}
	//
	//// if we got addrs then we'll update
	//if len(cAddrs) > 0 {
	//	config.Endpoints = cAddrs
	//}
	//
	//cli, err := clientv3.New(config)
	//if err != nil {
	//	return err
	//}
	//e.client = cli
	return nil
}

func encode(s *Service) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func decode(ds []byte) *Service {
	var s *Service
	json.Unmarshal(ds, &s)
	return s
}

func nodePath(s, id string) string {
	service := strings.Replace(s, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(prefix, service, node)
}

func servicePath(s string) string {
	return path.Join(prefix, strings.Replace(s, "/", "-", -1))
}

func (e *etcdv3Registry) Init(opts ...Option) error {
	return configure(e, opts...)
}

func (e *etcdv3Registry) Options() Options {
	return e.options
}

func (e *etcdv3Registry) Deregister(s *Service) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	e.Lock()
	// delete our hash of the service
	delete(e.register, s.Name)
	// delete our lease of the service
	delete(e.leases, s.Name)
	e.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	etcdClient := getEtcdClient()

	for _, node := range s.Nodes {
		_, err := etcdClient.Delete(ctx, nodePath(s.Name, node.Id))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *etcdv3Registry) Register(s *Service, opts ...RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	etcdClient := getEtcdClient()
	var leaseNotFound bool
	//refreshing lease if existing
	leaseID, ok := e.leases[s.Name]
	if ok {
		if _, err := etcdClient.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}

			// lease not found do register
			leaseNotFound = true
		}
	}

	// create hash of service; uint64
	h, err := hash.Hash(s, nil)
	if err != nil {
		return err
	}

	// get existing hash
	e.Lock()
	v, ok := e.register[s.Name]
	e.Unlock()

	// the service is unchanged, skip registering
	if ok && v == h && !leaseNotFound {
		return nil
	}

	service := &Service{
		Name:         s.Name,
		Version:      s.Version,
		Metadata:     s.Metadata,
		Endpoints:    s.Endpoints,
		RegistryTime: time.Now().Unix(),
		ProjectName:  s.ProjectName,
	}

	var options RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	if options.TTL.Seconds() > 0 {
		lgr, err = etcdClient.Grant(ctx, int64(options.TTL.Seconds()))
		if err != nil {
			return err
		}
	}

	for _, node := range s.Nodes {
		service.Nodes = []*Node{node}
		if lgr != nil {
			_, err = etcdClient.Put(ctx, nodePath(service.Name, node.Id), encode(service), clientv3.WithLease(lgr.ID))
		} else {
			_, err = etcdClient.Put(ctx, nodePath(service.Name, node.Id), encode(service))
		}
		if err != nil {
			return err
		}
	}

	e.Lock()
	// save our hash of the service
	e.register[s.Name] = h
	// save our leaseID of the service
	if lgr != nil {
		e.leases[s.Name] = lgr.ID
	}
	e.Unlock()

	return nil
}

func (e *etcdv3Registry) GetService(name string) ([]*Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	etcdClient := getEtcdClient()

	rsp, err := etcdClient.Get(ctx, servicePath(name)+"/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return nil, ErrNotFound
	}

	serviceMap := map[string]*Service{}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			s, ok := serviceMap[sn.Version]
			if !ok {
				s = &Service{
					Name:        sn.Name,
					Version:     sn.Version,
					Metadata:    sn.Metadata,
					Endpoints:   sn.Endpoints,
					ProjectName: sn.ProjectName,
				}
				serviceMap[s.Version] = s
			}

			for _, node := range sn.Nodes {
				s.Nodes = append(s.Nodes, node)
			}
		}
	}

	var services []*Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

// 原始版本不带任何信息，只带服务名称
// 更改成了带上Node的metaData 信息.
func (e *etcdv3Registry) ListServices() ([]*Service, error) {
	var services []*Service

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	etcdClient := getEtcdClient()

	rsp, err := etcdClient.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return []*Service{}, nil
	}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			service := &Service{
				Name:        sn.Name,
				Version:     sn.Version,
				Metadata:    sn.Metadata,
				Endpoints:   sn.Endpoints,
				Nodes:       sn.Nodes,
				ProjectName: sn.ProjectName,
			}
			services = append(services, service)
		}
	}

	//for k := range nameSet {
	//	service := &Service{}
	//	service.Name = k
	//	services = append(services, service)
	//}

	return services, nil
}

func (e *etcdv3Registry) Watch(opts ...WatchOption) (Watcher, error) {
	return newEtcdv3Watcher(e, e.options.Timeout, opts...)
}

func (e *etcdv3Registry) String() string {
	return "etcdv3"
}

func NewEtcdV3Registry(opts ...Option) Registry {
	e := &etcdv3Registry{
		options:  Options{},
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}
	configure(e, opts...)
	return e
}
