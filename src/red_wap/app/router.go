package app

import (
	controllers "baseGo/src/red_wap/app/controller"
	"baseGo/src/red_wap/app/middleware"
	"baseGo/src/red_wap/app/server"
	"baseGo/src/red_wap/conf"
	"time"
)

var (
	userCtrl = new(controllers.UserController)
	// loginCtrl = new(controllers.LoginController)

// roomCtrl             = new(controllers.RoomController)
// chatCtrl             = new(controllers.ChatController)
// msgCtrl              = new(controllers.MessageHistoryController)
// redPacketCtl         = new(controllers.RedPacketController)
// postCtl              = new(controllers.PostController)
// acCtl                = new(controllers.ActivePictureController)
// MemberMessageCtl     = new(controllers.MemberMessageControllers)
// UploadFileController = new(controllers.UploadFileController)
)

func RegisteRouter(echo *server.Echo) {

	echo.GET("/version", func(ctx server.Context) error {
		return ctx.JSON(200, map[string]interface{}{"version": conf.Version, "time": time.Now().UnixNano()})
	})
	// echo.Any("/*", UploadFileController.DownLoadFile) // 文件下载

	// noAuthApiPc := echo.Group("/api/wap", middleware.TimeLog, middleware.NotAuthInit)
	// {
	// 	noAuthApiPc.POST("/login", loginCtrl.login)
	// // 机器人初始化
	// noAuthApiPc.POST("/checkOrder", redPacketCtl.CheckOrder)
	// // 异常注单结算
	// noAuthApiPc.POST("/settleOrder", redPacketCtl.SettleOrder)
	// // 注单统计数据处理
	// noAuthApiPc.POST("/statisticOrder",redPacketCtl.StatisticalOrder)
	// }

	WapAuthApiApp := echo.Group("/api/wap", middleware.TimeLog, middleware.WapAuthInit)
	{
		// WapAuthApiApp.POST("/im/conn", userCtrl.ImCoon)                   // 获取im连接数据
		// 	WapAuthApiApp.POST("/im/upload", UploadFileController.UpLoadFile) // 在im上传文件
		// 	// 聊天
		// 	WapAuthApiApp.POST("/chat/joinroom", chatCtrl.ChangeRoom)                   // 进入房间
		// 	WapAuthApiApp.POST("/chat/outroom", chatCtrl.OutRoom)                       // 退出房间
		// 	WapAuthApiApp.POST("/chat/sendmessage", chatCtrl.SendMessage)               // 发消息
		// 	WapAuthApiApp.POST("/chat/messagehistory", msgCtrl.GetMessageHistory)       // 历史消息
		// 	WapAuthApiApp.POST("/chat/joinRoomByGroupId", chatCtrl.ChangeRoomByGroupId) // 进入私密房间

		WapAuthApiApp.POST("/info", userCtrl.GetUserInfo) // 获取会员信息

		// 	// 红包群列表
		// 	WapAuthApiApp.POST("/room/list", roomCtrl.RoomList)               // 群列表
		// 	WapAuthApiApp.POST("/agency/queryRoomOne", roomCtrl.QueryRoomOne) // 查询单个群信息

		// 	// 红包
		// 	WapAuthApiApp.POST("/red/add", redPacketCtl.CreateRedPacket)     // 发红包
		// 	WapAuthApiApp.POST("/red/receive", redPacketCtl.GrabRedEnvelope) // 抢红包
		// 	WapAuthApiApp.POST("/red/list", redPacketCtl.GetRedList)         // 红包记录
		// 	WapAuthApiApp.POST("/red/info", redPacketCtl.GetRedInfo)         // 红包详情
		// 	//WapAuthApiApp.POST("/red/js", redPacketCtl.GetRedjs)             // 红包结算

		// 	// 公告
		// 	WapAuthApiApp.POST("/post/list", postCtl.GetPostList)                      // 公告列表
		// 	WapAuthApiApp.POST("/active/list", acCtl.GetActiveList)                    // 活动记录列表
		// 	WapAuthApiApp.POST("/message/list", MemberMessageCtl.GetMemberMessageList) // 历史消息列表
	}
}
