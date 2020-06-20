package webserver

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/modules"
	"baseGo/src/red_agency/app"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	config2 "baseGo/src/red_agency/conf"
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
