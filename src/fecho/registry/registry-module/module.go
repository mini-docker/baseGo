package registry_module

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"fecho/golog"
	"fecho/modules"
	"fecho/registry"
	"fecho/utility/uuid"
)

var (
	rm            *registryModule
	LogicGrpcAddr = make(map[string]string)
	CometGrpcAddr = make(map[string]string)
)

var EtcdConn = struct {
	sync.RWMutex
	m map[string]map[string]string
}{m: make(map[string]map[string]string)}

type registryModule struct {
	m           *modules.Module
	r           registry.Registry
	registryTTL time.Duration
	w           registry.Watcher

	service *registry.Service
	sync.RWMutex

	stop chan bool
	tk   *time.Ticker
}

func init() {
	rm = new(registryModule)
	rm.stop = make(chan bool, 1)
	rm.m = modules.Register("registry", 37)
	rm.registryTTL = time.Second * 5
}

func Start(opt ...Option) {

	opts := Options{}

	for _, o := range opt {
		o(&opts)
	}
	// first time init registry
	rm.service = &registry.Service{
		Name:    opts.containerId,
		Version: opts.version,
		Metadata: map[string]string{
			"sites":    opts.sites,
			"addr":     "http://" + fmt.Sprintf("%s:%d", opts.addr, opts.port),
			"tcpAddr":  opts.tcpAddr,
			"wsAddr":   opts.wsAddr,
			"httpAddr": opts.httpAddr,
		},
		Nodes: []*registry.Node{
			{
				Id:      uuid.NewV4().String(),
				Address: "http://" + opts.addr,
				Port:    opts.port,
			},
		},
		RegistryTime: time.Now().Unix(),
		ProjectName:  opts.name,
	}

	if opts.ttl != 0 {
		rm.registryTTL = time.Second * time.Duration(opts.ttl)
	}

	rm.r = registry.NewRegistry(
		registry.Addrs(opts.registryAddr...),
		registry.Timeout(time.Duration(opts.ttl*2)*time.Second), // default 2 times than ttl
	)

	rm.tk = time.NewTicker(rm.registryTTL)

	err := rm.r.Register(rm.service, registry.RegisterTTL(rm.registryTTL))

	if err != nil {
		golog.Error("registry", "Start", "registry error:%+v", err)
		rm.m.StopComplete()
		return
	}

	go rm.keepAlive()
	go rm.run()
	go rm.watch()

}

func (rm *registryModule) keepAlive() {
	for {
		select {
		case <-rm.tk.C:
			rm.Lock()
			serv := *rm.service // copy
			rm.Unlock()
			err := rm.r.Register(&serv, registry.RegisterTTL(rm.registryTTL))
			if err != nil {
				golog.Error("registry", "Start", "registry error:%+v", err)
				return
			}
		case <-rm.stop:
			rm.tk.Stop()
			return
		}
	}
}

func (rm *registryModule) run() {
	<-rm.m.Stop
	rm.stop <- true

	rm.Lock()
	serv := *rm.service
	rm.Unlock()

	rm.r.Deregister(&serv)
	rm.m.StopComplete()
}

func (rm *registryModule) watch() {
	watcher, err := registry.Watch()
	if err != nil {
		panic(err)
	}
	for {
		result, err := watcher.Next()
		if err != nil {
			panic(err)
		}
		EtcdConn.RLock()
		switch result.Action {
		case "create":
			EtcdConn.m[result.Service.ProjectName] = result.Service.Metadata
			if strings.Contains(result.Service.ProjectName, "red-lgc") {
				LogicGrpcAddr[result.Service.ProjectName] = strings.Replace(result.Service.Metadata["addr"], "http://", "", -1)
			}
			if strings.Contains(result.Service.ProjectName, "red-met") {
				CometGrpcAddr[result.Service.ProjectName] = strings.Replace(result.Service.Metadata["addr"], "http://", "", -1)
			}
		case "update":
			// 更新缓存
			EtcdConn.m[result.Service.ProjectName] = result.Service.Metadata
			if strings.Contains(result.Service.ProjectName, "red-lgc") {
				if LogicGrpcAddr[result.Service.ProjectName] != strings.Replace(result.Service.Metadata["addr"], "http://", "", -1) {
					LogicGrpcAddr[result.Service.ProjectName] = strings.Replace(result.Service.Metadata["addr"], "http://", "", -1)
				}
			}
			if strings.Contains(result.Service.ProjectName, "red-met") {
				if CometGrpcAddr[result.Service.ProjectName] != strings.Replace(result.Service.Metadata["addr"], "http://", "", -1) {
					CometGrpcAddr[result.Service.ProjectName] = strings.Replace(result.Service.Metadata["addr"], "http://", "", -1)
				}
			}
		case "delete":
			// 移除缓存
			delete(EtcdConn.m, result.Service.ProjectName)
			if strings.Contains(result.Service.ProjectName, "red-lgc") {
				delete(LogicGrpcAddr, result.Service.ProjectName)
			}
			if strings.Contains(result.Service.ProjectName, "red-met") {
				delete(CometGrpcAddr, result.Service.ProjectName)
			}
		default:

		}
		EtcdConn.RUnlock()
	}
}

func GetLogicHttpUrl() string {
	EtcdConn.RLock()
	defer EtcdConn.RUnlock()
	etcdConf := EtcdConn.m
	if len(etcdConf) > 0 {
		for k, v := range etcdConf {
			if strings.Contains(k, "red-lgc") {
				return v["httpAddr"]
			}
		}
	}
	return ""
}

func GetCometTcpUrl() string {
	EtcdConn.RLock()
	defer EtcdConn.RUnlock()
	etcdConf := EtcdConn.m
	if len(etcdConf) > 0 {
		for k, v := range etcdConf {
			if strings.Contains(k, "red-met") {
				return v["tcpAddr"]
			}
		}
	}
	return ""
}

func GetCometRpcUrl() string {
	EtcdConn.RLock()
	defer EtcdConn.RUnlock()
	etcdConf := EtcdConn.m
	if len(etcdConf) > 0 {
		for k, v := range etcdConf {
			if strings.Contains(k, "red-met") {
				return strings.Replace(v["addr"], "http://", "", -1)
			}
		}
	}
	return ""
}

func GetCometWsUrl() string {
	EtcdConn.RLock()
	defer EtcdConn.RUnlock()
	etcdConf := EtcdConn.m
	if len(etcdConf) > 0 {
		for k, v := range etcdConf {
			if strings.Contains(k, "red-met") {
				return v["wsAddr"]
			}
		}
	}
	return ""
}
