package registry_module

import (
	"fmt"
	"strings"

	"github.com/mini-docker/baseGo/fecho/golog"
	"github.com/mini-docker/baseGo/fecho/registry"
)

type Service struct {
	HostPort string
}

func GetSelf() registry.Service {
	rm.RLock()
	srv := *rm.service
	rm.RUnlock()
	return srv
}

// 获取站点对应的服务列表.
func GetServices() ([]Service, error) {
	var res = make([]Service, 0)
	services, err := registry.ListServices()
	if err != nil {
		golog.Error("module", "GetNodeList", "get Nodelist failed :%v", err)
		return res, err
	}

	for _, service := range services {
		s := Service{
			HostPort: service.Metadata["addr"],
		}
		res = append(res, s)
	}
	return res, nil
}

// 获取不含平台的所有节点.
func GetNodeList() map[string]string {
	var res = make(map[string]string)

	services, err := registry.ListServices()
	if err != nil {
		golog.Error("module", "GetNodeList", "get Nodelist failed :%v", err)
		return res
	}

	for _, service := range services {
		if strings.Contains(service.Metadata["sites"], "0000") ||
			strings.Contains(service.ProjectName, "job") ||
			strings.Contains(service.ProjectName, "logic") ||
			strings.Contains(service.ProjectName, "comet") {
			continue
		}
		res[service.ProjectName] = service.Metadata["addr"] + "-" + service.Metadata["sites"]
	}
	return res
}

// 获取所有服务  [0000]:Service
func ListServices() map[string][]Service {

	var res = make(map[string][]Service)
	services, err := registry.ListServices()

	if err != nil {
		golog.Error("module", "GetNodeList", "get Nodelist failed :%v", err)
		return res
	}

	for _, service := range services {
		ids := strings.Split(service.Metadata["sites"], ",")
		if len(ids) == 0 {
			continue
		}
		node := service.Nodes[0]
		for _, v := range ids {
			if len(v) != 4 {
				continue
			}
			s := Service{
				HostPort: fmt.Sprintf("%s:%d", node.Address, node.Port),
			}

			if _, ok := res[v]; !ok {
				res[v] = make([]Service, 0)
			}
			var has = false
			// 如果发现已经存在. 就不存储了.
			for _, v := range res[v] {
				if v.HostPort == s.HostPort {
					has = true
					continue
				}
			}
			if !has {
				res[v] = append(res[v], s)
			}
		}
	}
	return res
}

// 动态 更新sites
func ReloadSites(sites string) {
	rm.Lock()
	defer rm.Unlock()
	rm.service.Metadata["sites"] = sites
}
