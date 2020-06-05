package webserver

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/modules"
	"baseGo/src/red_admin/app"
	"baseGo/src/red_admin/app/middleware/validate"
	"baseGo/src/red_admin/app/server"
	"net"

	config2 "baseGo/src/red_admin/conf"
)

var (
	webServerModule *modules.Module
	webServer       *echo.Echo
	version         string
	httpPort        string
)

// http tcp rpc proto https协议
func Start(port string) {
	webServerModule = modules.Register("WebServer", 3)
	httpPort = port

	go run()

	<-webServerModule.Stop
	webServerModule.StopComplete()
}

func run() {
	e := server.New()
	if config2.GetAppConfig().AppEnvironment == "dev" {
		e.Debug = true
	}
	e.Validator = &validate.FrontValidate{}

	addr, err := net.ResolveTCPAddr("tcp", httpPort)
	if err != nil {
		panic("api" + "run" + "invalid config address " + httpPort + err.Error())
	}
	app.RegisteRouter(e)
	//启动服务
	err = e.Start(addr.String())
	if err != nil {
		panic("api" + "run" + "NewTemplateRenderer error " + httpPort + err.Error())
	}
}
