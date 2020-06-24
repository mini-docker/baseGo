package app

import (
	controllers "baseGo/src/red_agency/app/controller"
	"baseGo/src/red_agency/app/middleware"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/conf"
	"time"
)

var (
	userCtrl             = new(controllers.UserController)
	roomCtrl             = new(controllers.RoomController)
	postCtl              = new(controllers.PostController)
	acCtl                = new(controllers.ActivePictureController)
	loginCtrl            = new(controllers.LoginController)
	agencyCtrl           = new(controllers.AgencyController)
	RedOrderRecordCtl    = new(controllers.RedOrderRecordController)
	UploadFileController = new(controllers.UploadFileController)
	OrdinaryRedPacketCtl = new(controllers.OrdinaryRedPacketController)
	RobotCtl             = new(controllers.RobotController)
	RedPacketSiteCtl     = new(controllers.RedPacketSiteController)
	NewSiteCtl           = new(controllers.NewSiteController)
	LogCtl               = new(controllers.LogController)
	OrderStatistical     = new(controllers.OrderStatistical)
)

func RegisteRouter(echo *server.Echo) {

	echo.GET("/version", func(ctx server.Context) error {
		return ctx.JSON(200, map[string]interface{}{"version": conf.Version, "time": time.Now().UnixNano()})
	})
	echo.Any("/*", UploadFileController.DownLoadFile) // 文件下载
	// pc api路由 非验证
	noAuthApiPc := echo.Group("/api", middleware.TimeLog, middleware.NotAuthInit)
	{
		noAuthApiPc.POST("/agency/login", loginCtrl.Login)
		noAuthApiPc.POST("/agency/orders", OrdinaryRedPacketCtl.Orders)
		// 数据初始化
		noAuthApiPc.POST("/agency/newSite", NewSiteCtl.NewSite)
		// 机器人初始化
		noAuthApiPc.POST("/agency/newRobot", NewSiteCtl.NewRobot)
		// 统计数据初始化
		noAuthApiPc.POST("/agency/initStatistical", NewSiteCtl.InitStatistical)
	}
	AuthApiApp := echo.Group("/api", middleware.TimeLog, middleware.AgencyAuthInit)
	{
		AuthApiApp.POST("/agency/logout", loginCtrl.Logout) // 注销

		// 游戏群管理ok
		AuthApiApp.POST("/agency/queryRoomList", roomCtrl.QueryRoomList)   // 游戏群列表ok
		AuthApiApp.POST("/agency/addRoom", roomCtrl.AddRoom)               // 添加群
		AuthApiApp.POST("/agency/queryRoomOne", roomCtrl.QueryRoomOne)     // 查询单个群信息
		AuthApiApp.POST("/agency/editRoom", roomCtrl.EditRoom)             // 修改群信息
		AuthApiApp.POST("/agency/editRoomStatus", roomCtrl.EditRoomStatus) // 修改群状态
		AuthApiApp.POST("/agency/delRoom", roomCtrl.DelRoom)               // 删除群
		AuthApiApp.POST("/agency/addRed", roomCtrl.AddRed)                 // 发送普通红包
		AuthApiApp.POST("/agency/roomCode", roomCtrl.RoomCode)             // 群枚举

		// 代理管理ok
		AuthApiApp.POST("/agency/queryAgencyList", agencyCtrl.QueryAgencyList)         // 查询代理列表ok
		AuthApiApp.POST("/agency/addAgency", agencyCtrl.AddAgency)                     // 添加代理
		AuthApiApp.POST("/agency/queryAgencyOne", agencyCtrl.QueryAgencyOne)           // 查询单个代理信息
		AuthApiApp.POST("/agency/editAgency", agencyCtrl.EditAgency)                   // 修改代理信息
		AuthApiApp.POST("/agency/editAgencyStatus", agencyCtrl.EditAgencyStatus)       // 修改代理状态
		AuthApiApp.POST("/agency/resetAgencyPassword", agencyCtrl.ResetAgencyPassword) // 重置代理密码
		AuthApiApp.POST("/agency/delAgency", agencyCtrl.DelAgency)                     // 删除代理
		AuthApiApp.POST("/agency/editPassword", agencyCtrl.EditPassword)               // 修改密码

		// 会员管理ok
		AuthApiApp.POST("/agency/queryUserList", userCtrl.QueryUserList)     // 查询会员列表 ok
		AuthApiApp.POST("/agency/kickUsers", userCtrl.KickUsers)             // 批量踢线
		AuthApiApp.POST("/agency/editUsersStatus", userCtrl.EditUsersStatus) // 批量启用/停用会员

		// 注单管理
		AuthApiApp.POST("/agency/queryRedOrderRecordList", RedOrderRecordCtl.QueryRedRecordList) // 注单列表 ok
		AuthApiApp.POST("/agency/redInfo", RedOrderRecordCtl.GetRedInfo)                         // 红包详情

		// 公告管理
		AuthApiApp.POST("/agency/queryPostList", postCtl.GetAgencyPostList) // 公告列表
		AuthApiApp.POST("/agency/addPost", postCtl.AddPost)                 // 添加公告
		AuthApiApp.POST("/agency/queryPostById", postCtl.QueryPostById)     // 根据id查询公告
		AuthApiApp.POST("/agency/editPost", postCtl.EditPost)               // 修改公告
		AuthApiApp.POST("/agency/editPostStatus", postCtl.EditPostStatus)   // 修改公告状态
		AuthApiApp.POST("/agency/delPost", postCtl.DelPost)                 // 删除公告

		// 活动/广告管理
		AuthApiApp.POST("/agency/queryActiveList", acCtl.GetAgencyActiveList) // 活动列表
		AuthApiApp.POST("/agency/addActive", acCtl.AddActive)                 // 添加活动
		AuthApiApp.POST("/agency/queryActiveById", acCtl.QueryActiveById)     // 根据id查询活动
		AuthApiApp.POST("/agency/editActive", acCtl.EditActive)               // 修改活动
		AuthApiApp.POST("/agency/editActiveStatus", acCtl.EditActiveStatus)   // 修改活动状态
		AuthApiApp.POST("/agency/delActive", acCtl.DelActive)                 // 删除活动

		AuthApiApp.POST("/agency/uploads", UploadFileController.UpLoadFile) // 在im上传文件
		AuthApiApp.POST("/agency/upload", UploadFileController.UpLoadFiles) // 在本地上传文件

		// 普通红包管理
		AuthApiApp.POST("/agency/ordinaryRed/list", OrdinaryRedPacketCtl.GetRedList)  // 查询红包列表
		AuthApiApp.POST("/agency/ordinaryRed/edit", OrdinaryRedPacketCtl.EditRedInfo) // 修改红包
		AuthApiApp.POST("/agency/ordinaryRed/info", OrdinaryRedPacketCtl.GetRedInfo)  // 查看红包领取情况
		AuthApiApp.POST("/agency/ordinaryRed/del", OrdinaryRedPacketCtl.DelRedInfo)   // 删除红包

		// 机器人管理
		AuthApiApp.POST("/agency/robot/list", RobotCtl.QueryRobotList)                 // 查询机器人列表
		AuthApiApp.POST("/agency/robot/getRobotAccounts", RobotCtl.CreatRobotAccounts) // 批量生成机器人账号
		AuthApiApp.POST("/agency/robot/insertRobots", RobotCtl.InsertRobots)           // 批量生成机器人
		AuthApiApp.POST("/agency/robot/delRobots", RobotCtl.DelRobots)                 // 批量删除机器人

		// 站点管理
		AuthApiApp.POST("/agency/site/list", RedPacketSiteCtl.QueryPacketSiteList)               // 站点列表
		AuthApiApp.POST("/agency/site/addSite", RedPacketSiteCtl.AddPacketSite)                  // 添加站点
		AuthApiApp.POST("/agency/site/editSite", RedPacketSiteCtl.EditPacketSite)                // 修改站点
		AuthApiApp.POST("/agency/site/editSiteStatus", RedPacketSiteCtl.EditRedPacketSiteStatus) // 修改站点状态
		AuthApiApp.POST("/agency/site/delSite", RedPacketSiteCtl.DelPacketSite)                  // 删除站点
		AuthApiApp.POST("/agency/agencyCode", RedPacketSiteCtl.SiteCode)                         // 站点枚举

		// 操作日志
		AuthApiApp.POST("/agency/log/list", LogCtl.QueryLogs) // 操作日志列表

		// 统计
		AuthApiApp.POST("/agency/orderStatistical", OrderStatistical.QueryOrderStatistical) // 盈利分析
	}
}
