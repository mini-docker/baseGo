package webserver

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/modules"
	"baseGo/src/red_robot/app/server"
	"net"
)

var (
	webServerModule *modules.Module
	webServer       *echo.Echo
	version         string
	httpPort        string
)

func Start(port string) {
	webServerModule = modules.Register("WebServer", 3)
	httpPort = port

	go run()

	<-webServerModule.Stop
	webServerModule.StopComplete()
}

func run() {
	e := server.New()
	addr, err := net.ResolveTCPAddr("tcp", httpPort)
	if err != nil {
		panic("api" + "run" + "invalid config address " + httpPort + err.Error())
	}
	//启动服务
	err = e.Start(addr.String())
	if err != nil {
		panic("api" + "run" + "NewTemplateRenderer error " + httpPort + err.Error())
	}
}
