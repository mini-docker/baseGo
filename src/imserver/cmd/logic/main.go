package main

import (
	"baseGo/src/fecho/cli"
	modules2 "baseGo/src/fecho/modules"
	"baseGo/src/fecho/registry"
	"baseGo/src/imserver/pkg/etcd"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	//"github.com/bilibili/discovery/naming/grpc"
	log "baseGo/src/fecho/golog"
	module "baseGo/src/fecho/registry/registry-module"
	"baseGo/src/imserver/internal/logic"
	"baseGo/src/imserver/internal/logic/conf"
	"baseGo/src/imserver/internal/logic/grpc"
	"baseGo/src/imserver/internal/logic/http"
	"baseGo/src/imserver/pkg/ip"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ver     = "v1.0.2"
	appName = "im-logic"
	cliName = "im-logic"
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
		//todo 本地测试
		if c.String("config_file") == "" {
			//c.Set("config_file", "/Users/yiwang/go/src/red-packet/src/imserver/cmd/logic/logic_conf.yaml")
			//c.Set("config_file", "./logic_conf.yaml")
			c.Set("config_file", "/Users/tongjunchao/goproduct/src/red-packet/src/imserver/cmd/logic/logic_conf.yaml")
		}
		//配置文件初始化
		cfg, err := conf.ParseConfigFile(c.String("config_file"))
		if err != nil {
			return err
		}

		return conf.Init(cfg)
	}

	rootCmd.Action = startServer

	err := rootCmd.Run(os.Args)
	if err != nil {
		log.Error("im-comet", "main", "", err)
		os.Exit(1)
	}
}

func startServer(c *cli.Context) error {
	//log.Info( "", "main", "goim-logic [version: %s env: %+v] start", ver, conf.Conf.Env)
	// logic
	srv := logic.New(conf.Conf)
	//go srv.Consume()
	// 启动离线推送定时器
	//go dao.InitOffMessage()
	etcd.InitEtcd(
		registry.Addrs(strings.Split(strings.Trim(conf.Conf.Registry.Addr, ","), ",")...),
		registry.Timeout(time.Duration(conf.Conf.Registry.TTL*2)*time.Second),
	)
	go Start()
	// 定时更新连接缓存
	go logic.UpdateCometConn()
	// 初始化缓存
	//go logic.SaveCache()
	//cancel := register(dis, srv)
	rpcSrv := grpc.New(conf.Conf.RPCServer, srv)
	httpSrv := http.New(conf.Conf.HTTPServer, srv) // http server init

	// signal
	signalCh := make(chan os.Signal)
	signal.Notify(
		signalCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	select {
	case <-signalCh:
		modules2.InitiateFullShutdown()
		srv.Close()
		httpSrv.Close()
		rpcSrv.GracefulStop()
		//log.Info( "", "main", "goim-logic [version: %s] exit", ver)
		//log.Flush()
		return nil
	case <-modules2.GlobalShutdown:
		modules2.InitiateFullShutdown()
	}
	return errors.New("server has been shutdown")
}

func Start() {
	addr := ip.InternalIP()
	_, port, _ := net.SplitHostPort(":" + conf.Conf.RPCServer.Addr)
	env := conf.Conf.Env
	portNum, _ := strconv.Atoi(port)
	module.Start(
		module.WithRegistryAddr(strings.Split(strings.Trim(conf.Conf.Registry.Addr, ","), ",")...),
		module.WithAddr(addr),                  // 内网 ip ， 这个配置在 node 上面，主要是靠这个配置通讯. ip:addr 获取.
		module.WithContainerId(env.Host),       // 容器id  & 与内网ip一个意思，这个配置在根节点  hostname获取
		module.WithPort(portNum),               // 应用程序的端口
		module.WithTTL(conf.Conf.Registry.TTL), // 过期时间
		module.WithVersion("red-lgc"),          // 版本号
		module.WithHttpAddr(fmt.Sprintf("%v:%v", addr, conf.Conf.HTTPServer.Addr)), // ws地址
		module.WithName(fmt.Sprintf("red-lgc-%v", env.Host)),                       // 服务名称
	)
}
