package app

import (
	"baseGo/src/red_api/app/controller"
	"baseGo/src/red_api/app/middleware"
	"baseGo/src/red_api/app/server"
	"baseGo/src/red_api/conf"
	"time"
)

var ()

func RegisteRouter(echo *server.Echo) {

	echo.GET("/version", func(ctx server.Context) error {
		return ctx.JSON(200, map[string]interface{}{"version": conf.Version, "time": time.Now().UnixNano()})
	})
	UserCtrl := new(controller.UserController)
	// FinanceCtrl := new(controller.FinanceController)
	// CollectCtrl := new(controller.RedPacketCollectController)
	// pc api路由 非验证
	noAuthApiPc := echo.Group("/api", middleware.TimeLog, middleware.NotAuthInit)
	{
		//账户
		noAuthApiPc.POST("/login", UserCtrl.Login)
		noAuthApiPc.POST("/register", UserCtrl.Register)
		noAuthApiPc.POST("/userInfo", UserCtrl.GetUserInfo)

		//资金
		// noAuthApiPc.POST("/transferredIn", FinanceCtrl.TransferredIn)
		// noAuthApiPc.POST("/transferredOut", FinanceCtrl.TransferredOut)

		// //采集
		// noAuthApiPc.POST("/collect", CollectCtrl.Collect)
		// noAuthApiPc.POST("/collectByDate", CollectCtrl.CollectByDate)
	}

}
