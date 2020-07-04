package main

import (
	"baseGo/src/fecho/cli"
	modules2 "baseGo/src/fecho/modules"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	//resolver "github.com/bilibili/discovery/naming/grpc"

	log "baseGo/src/fecho/golog"
	"baseGo/src/fecho/registry"
	module "baseGo/src/fecho/registry/registry-module"
	"baseGo/src/imserver/internal/comet"
	"baseGo/src/imserver/internal/comet/conf"
	"baseGo/src/imserver/internal/comet/grpc"
	"baseGo/src/imserver/pkg/etcd"
	"baseGo/src/imserver/pkg/ip"
	"time"
)

const (
	ver     = "v1.0.2"
	appName = "im-comet-js135"
	cliName = "im-comet"
)

func main() {
	rootCmd := cli.NewApp()
	rootCmd.Name = appName
	rootCmd.HelpName = cliName
	rootCmd.Version = ver
	rootCmd.Usage = ""

	rootCmd.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "MCONDIF_CONFIG_FILE",
			Name:   "config_file",
			Usage:  "appconfig file path",
		},
	}

	rootCmd.Before = func(c *cli.Context) error {
		if c.String("config_file") == "" {
			//c.Set("config_file", "/Users/yiwang/go/src/red-packet/src/imserver/cmd/comet/comet_conf.yaml")
			c.Set("config_file", "./comet_conf.yaml")
			// c.Set("config_file", "/Users/tongjunchao/goproduct/src/red-packet/src/imserver/cmd/comet/comet_conf.yaml")
		}
		//配置文件初始化
		cfg, err := conf.ParseConfigFile(c.String("config_file"))
		if err != nil {
			return err
		}
		return conf.Init(cfg)
	}

	rootCmd.Action = func(ctx *cli.Context) error {

		rand.Seed(time.Now().UTC().UnixNano())
		runtime.GOMAXPROCS(runtime.NumCPU())
		println(conf.Conf.Debug)
		//log.Info( "", "main", "goim-comet [version: %s env: %+v] start", ver, conf.Conf.Env)

		etcd.InitEtcd(
			registry.Addrs(strings.Split(strings.Trim(conf.Conf.Registry.Addr, ","), ",")...),
			registry.Timeout(time.Duration(conf.Conf.Registry.TTL*2)*time.Second),
		)
		go Start()
		// 定时更新连接缓存
		go comet.UpdateLogicConn()
		// new comet server
		srv := comet.NewServer(conf.Conf)
		if err := comet.InitWhitelist(conf.Conf.Whitelist); err != nil {
			log.Error("im-comet", "main", "goim-comet err:", err)
			return err
		}

		if err := comet.InitWebsocket(srv, conf.Conf.Websocket.Bind, runtime.NumCPU()); err != nil {
			return err
		}
		// new grpc server
		rpcSrv := grpc.New(conf.Conf.RPCServer, srv)
		//cancel := register(dis, srv)
		// signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			s := <-c
			//log.Info( "", "main", "goim-comet get a signal %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				modules2.InitiateFullShutdown()
				rpcSrv.GracefulStop()
				srv.Close()
				//log.Info( "", "main", "goim-comet [version: %s] exit", ver)
				//log.Flush()
				return nil
			case syscall.SIGHUP:
				modules2.InitiateFullShutdown()
			default:
				return nil
			}
		}
	}

	err := rootCmd.Run(os.Args)
	if err != nil {
		log.Error("im-comet", "main", "", err)
		os.Exit(1)
	}
}

func Start() {

	addr := ip.InternalIP()
	_, port, _ := net.SplitHostPort(":" + conf.Conf.RPCServer.Addr)
	env := conf.Conf.Env

	portNum, _ := strconv.Atoi(port)
	module.Start(
		module.WithRegistryAddr(strings.Split(strings.Trim(conf.Conf.Registry.Addr, ","), ",")...),
		module.WithAddr(addr),                                // 内网 ip ， 这个配置在 node 上面，主要是靠这个配置通讯. ip:addr 获取.
		module.WithContainerId(env.Host),                     // 容器id  & 与内网ip一个意思，这个配置在根节点  hostname获取
		module.WithPort(portNum),                             // 应用程序的端口
		module.WithTTL(conf.Conf.Registry.TTL),               // 过期时间
		module.WithVersion("red-met"),                        // 版本号
		module.WithWsAddr(conf.Conf.Websocket.PubUrl),        // ws地址
		module.WithTcpAddr(conf.Conf.TCP.PubUrl),             // tcp地址
		module.WithName(fmt.Sprintf("red-met-%v", env.Host)), // 服务名称
	)
}
