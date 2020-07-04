package main

import (
	"baseGo/src/fecho/cli"
	modules2 "baseGo/src/fecho/modules"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	//"github.com/bilibili/discovery/naming"
	"baseGo/src/imserver/internal/job"
	"baseGo/src/imserver/internal/job/conf"

	//resolver "github.com/bilibili/discovery/naming/grpc"
	log "baseGo/src/fecho/golog"
	"baseGo/src/fecho/registry"
	module "baseGo/src/fecho/registry/registry-module"
	"baseGo/src/imserver/pkg/etcd"
	"baseGo/src/imserver/pkg/ip"
	"strings"
	"time"
)

var (
	ver     = "v1.0.2"
	appName = "im-job"
	cliName = "im-job"
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
			//c.Set("config_file", "/Users/yiwang/go/src/red-packet/src/imserver/cmd/job/job_conf.yaml")
			c.Set("config_file", "./job_conf.yaml")
			// c.Set("config_file", "/Users/tongjunchao/goproduct/src/red-packet/src/imserver/cmd/job/job_conf.yaml")
		}
		//配置文件初始化
		cfg, err := conf.ParseConfigFile(c.String("config_file"))
		if err != nil {
			return err
		}

		return conf.Init(cfg)
	}

	rootCmd.Action = func(ctx *cli.Context) error {

		//log.Info( "", "main", "goim-job [version: %s env: %+v] start", ver, conf.Conf.Env)
		etcd.InitEtcd(
			registry.Addrs(strings.Split(strings.Trim(conf.Conf.Registry.Addr, ","), ",")...),
			registry.Timeout(time.Duration(conf.Conf.Registry.TTL*2)*time.Second),
		)
		go Start()

		// job
		j := job.New(conf.Conf)
		go j.Consume()
		// signal
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		for {
			s := <-c
			//log.Info( "im-comet", "main", "goim-job get a signal %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
				j.Close()
				modules2.InitiateFullShutdown()
				//etcdService.Stop()
				//log.Info( "im-comet", "main", "goim-job [version: %s] exit", ver)
				return nil
			case syscall.SIGHUP:
				modules2.InitiateFullShutdown()
				//etcdService.Stop()
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
	env := conf.Conf.Env
	module.Start(
		module.WithRegistryAddr(strings.Split(strings.Trim(conf.Conf.Registry.Addr, ","), ",")...),
		module.WithAddr(addr),                                // 内网 ip ， 这个配置在 node 上面，主要是靠这个配置通讯. ip:addr 获取.
		module.WithContainerId(env.Host),                     // 容器id  & 与内网ip一个意思，这个配置在根节点  hostname获取
		module.WithTTL(conf.Conf.Registry.TTL),               // 过期时间
		module.WithVersion("reb-job"),                        // 版本号
		module.WithName(fmt.Sprintf("reb-job-%v", env.Host)), // 服务名称
	)
}
