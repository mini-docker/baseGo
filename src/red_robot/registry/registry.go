package registry

import (
	"baseGo/src/red_robot/conf"
	"fmt"
	"net"
	"os"
	"strings"

	"baseGo/src/fecho/golog"
	module "baseGo/src/fecho/registry/registry-module"

	"github.com/pkg/errors"
)

func Start() {
	//TODO ？？
	ip := getIntranetIp()
	if ip == "" {
		golog.Error("registry", "Start", "%v", errors.New("get ip failed"))
		return
	}
	host := getHostName()
	if host == "" {
		host = ip
	}
	module.Start(
		module.WithRegistryAddr(strings.Split(strings.Trim(conf.GetRegistryConfig().Addr, ","), ",")...),
		module.WithAddr(ip),                                // 内网 ip ， 这个配置在 node 上面，主要是靠这个配置通讯. ip:addr 获取.
		module.WithContainerId(host),                       // 容器id  & 与内网ip一个意思，这个配置在根节点  hostname获取
		module.WithPort(conf.GetAppConfig().ApiPort),       // 应用程序的端口
		module.WithTTL(conf.GetRegistryConfig().TTL),       // 过期时间
		module.WithVersion("red-robot"),                    // 版本号
		module.WithName(fmt.Sprintf("red-robot-%v", host)), // 服务名称
	)
}

func getIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getHostName() string {
	host, _ := os.Hostname()
	return host
}
