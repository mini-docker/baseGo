package main

import (
	"baseGo/fecho/cli"
	"baseGo/fecho/golog"
	"baseGo/fecho/modules"
	"baseGo/fecho/utility"
	"baseGo/red_admin/conf"
	"baseGo/red_admin/webserver"
	"errors"
	"os"
	"os/signal"

	// "red_admin/app/middleware"
	// "red_admin/registry"
	"runtime/pprof"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ver     = "v1.0.2"
	appName = "red-admin"
	cliName = "red-admin"
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
			//c.Set("config_file", "/Users/tongjunchao/goproduct/src/red-packet/src/red_admin/red_admin_conf.yaml")
			c.Set("config_file", "./red_admin_conf.yaml")
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
		golog.Error("chat-server", "main", "", err)
		os.Exit(1)
	}
}

func startServer(c *cli.Context) error {
	//注册日志
	golog.Logger.From = cliName
	modules.RegisterLogger(golog.Logger)

	//api接口启动
	go webserver.Start(conf.GetAppConfig().Addr + ":" + utility.ToStr(conf.GetAppConfig().ApiPort))

	//服务注册
	// go registry.Start()

	// 在线状态检测
	// go middleware.InitOnlineCheck()

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
		//golog.Debug("main", "server", "program was interrupted, shutting down.")
		modules.InitiateFullShutdown()
	case <-modules.GlobalShutdown:
	}

	// wait for shutdown to complete, panic after timeout
	time.Sleep(5 * time.Second)
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)

	return errors.New("server has been shutdown")
}