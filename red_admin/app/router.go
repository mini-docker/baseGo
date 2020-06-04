package app

import (
	"baseGo/red_admin/app/controller"
	"baseGo/red_admin/app/middleware"
	"baseGo/red_admin/app/server"
	"baseGo/red_admin/conf"
	"time"
)

var (
	// SystemMenuController          = new(controller.SystemMenuController)
	// SystemRoleController          = new(controller.SystemRoleController)
	SystemAdminController = new(controller.SystemAdminController)

// SystemLineMealController      = new(controller.SystemLineMealController)
// SystemLineController          = new(controller.SystemLineController)
// -SystemAgencyController = new(controller.SystemAgencyController)
// SystemGameController          = new(controller.SystemGameController)
// SystemRoomController          = new(controller.SystemRoomController)
// SystemUserController          = new(controller.SystemUserController)
// RedOrderRecordController      = new(controller.RedOrderRecordController)
// SystemLineRoyaltyController   = new(controller.SystemLineRoyaltyController)
// SystemPostController          = new(controller.SystemPostController)
// SystemActivePictureController = new(controller.SystemActivePictureController)
// OrderStatistical              = new(controller.OrderStatistical)
)

func RegisteRouter(echo *server.Echo) {

	echo.GET("/version", func(ctx server.Context) error {
		return ctx.JSON(200, map[string]interface{}{"version": conf.Version, "time": time.Now().UnixNano()})
	})

	// pc api路由 非验证
	// 登录后发sessiona给前端
	noAuthApiPc := echo.Group("/api", middleware.TimeLog, middleware.NotAuthInit)
	{
		noAuthApiPc.POST("/system/login", SystemAdminController.Login) // 登陆
	}
	// 之后再一次根据session 验证登录信息
	AuthApiApp := echo.Group("/api", middleware.TimeLog, middleware.AuthInit)
	{
		// 系统菜单管理ok
		// AuthApiApp.POST("/system/queryMenuList", SystemMenuController.QueryMenuList)         // 查询菜单列表
		// AuthApiApp.POST("/system/addMenu", SystemMenuController.AddMenu)                     // 添加菜单
		// AuthApiApp.POST("/system/queryMenuOne", SystemMenuController.QueryMenuOne)           // 查询单个菜单
		// AuthApiApp.POST("/system/queryChildrenById", SystemMenuController.QueryChildrenById) // 根据id查询子菜单
		// AuthApiApp.POST("/system/editMenu", SystemMenuController.EidtMenu)                   // 修改菜单

		// 角色管理ok
		// AuthApiApp.POST("/system/queryRoleList", SystemRoleController.QueryRoleList)             // 查询角色列表ok
		// AuthApiApp.POST("/system/addRole", SystemRoleController.AddRole)                         // 添加角色
		// AuthApiApp.POST("/system/queryRoleOne", SystemRoleController.QueryRoleOne)               // 查询单个角色
		// AuthApiApp.POST("/system/editRole", SystemRoleController.EidtRole)                       // 修改角色
		// AuthApiApp.POST("/system/editRoleStatus", SystemRoleController.EidtRoleStatus)           // 修改角色状态
		// AuthApiApp.POST("/system/queryRolePermission", SystemRoleController.QueryRolePermission) // 获取角色权限
		// AuthApiApp.POST("/system/setRolePermission", SystemRoleController.SetRolePermission)     // 角色赋权
		// AuthApiApp.POST("/system/delRole", SystemRoleController.DelRole)                         // 删除角色
		// AuthApiApp.POST("/system/queryRoleCode", SystemRoleController.QuerySystemRoleCode)       // 角色枚举

		// 系统账号管理ok
		// AuthApiApp.POST("/system/queryAdminList", SystemAdminController.QueryAdminList) // 查询系统管理员列表ok
		// AuthApiApp.POST("/system/addAdmin", SystemAdminController.AddAdmin)             // 添加系统管理员
		// AuthApiApp.POST("/system/queryAdminOne", SystemAdminController.QueryAdminOne)   // 查询单个系统管理员
		// AuthApiApp.POST("/system/editAdmin", SystemAdminController.EidtAdmin)           // 修改系统管理员
		// AuthApiApp.POST("/system/resetPassword", SystemAdminController.ResetPassword)   // 重置密码
		// AuthApiApp.POST("/system/delAdmin", SystemAdminController.DelAdmin)             // 删除管理员
		// AuthApiApp.POST("/system/editPassword", SystemAdminController.EditPassword)     // 修改密码

		// 登陆注销模块ok
		// AuthApiApp.POST("/system/owner", SystemRoleController.QueryRoleMenu) // 获取菜单权限
		// AuthApiApp.POST("/system/logout", SystemAdminController.Logout)      // 注销

		// 线路套餐管理ok
		// AuthApiApp.POST("/system/queryLineMealList", SystemLineMealController.QueryLineMealList)    // 查询线路套餐列表ok
		// AuthApiApp.POST("/system/addLineMeal", SystemLineMealController.AddLineMeal)                // 添加线路套餐
		// AuthApiApp.POST("/system/queryLineMealOne", SystemLineMealController.QueryLineMealOne)      // 查询单个线路套餐
		// AuthApiApp.POST("/system/editLineMeal", SystemLineMealController.EidtLineMeal)              // 修改线路套餐
		// AuthApiApp.POST("/system/queryLineMealCode", SystemLineMealController.QueryAllLineMealCode) // 获取套餐枚举

		// 线路管理ok
		// AuthApiApp.POST("/system/queryLineList", SystemLineController.QueryLineList)   // 查询线路列表ok
		// AuthApiApp.POST("/system/queryLineCode", SystemLineController.QueryAllLineId)  // 获取线路id枚举
		// AuthApiApp.POST("/system/addLine", SystemLineController.AddLine)               // 添加线路
		// AuthApiApp.POST("/system/queryLineOne", SystemLineController.QueryLineOne)     // 查询单个线路
		// AuthApiApp.POST("/system/editLine", SystemLineController.EidtLine)             // 修改线路
		// AuthApiApp.POST("/system/editLineStatus", SystemLineController.EidtLineStatus) // 修改线路状态

		// 业主超管管理ok
		// AuthApiApp.POST("/system/queryAgencyAdminList", SystemAgencyController.QueryAgencyAdminList) // 查询超管列表ok
		// AuthApiApp.POST("/system/addAgencyAdmin", SystemAgencyController.AddAgencyAdmin)             // 添加超管
		// AuthApiApp.POST("/system/queryAgencyAdminOne", SystemAgencyController.QueryAgencyAdminOne)   // 查询单个超管
		// AuthApiApp.POST("/system/editAgencyAdmin", SystemAgencyController.EditAgencyAdmin)           // 修改超管

		// 业主代理管理ok
		// AuthApiApp.POST("/system/QueryAgencyList", SystemAgencyController.QueryAgencyList)   // 查询代理列表ok
		// AuthApiApp.POST("/system/editAgencyStatus", SystemAgencyController.EditAgencyStatus) // 修改代理状态
		// AuthApiApp.POST("/system/agencyCode", SystemAgencyController.SiteCode)               // 站点枚举

		// 会员管理
		// AuthApiApp.POST("/system/queryUserList", SystemUserController.QueryUserList)     // 查询会员列表 ok
		// AuthApiApp.POST("/system/kickUsers", SystemUserController.KickUsers)             // 批量踢线
		// AuthApiApp.POST("/system/editUsersStatus", SystemUserController.EditUsersStatus) // 批量启用/停用会员

		// 游戏管理ok
		// AuthApiApp.POST("/system/queryGameList", SystemGameController.QueryGameList)   // 查询游戏列表ok
		// AuthApiApp.POST("/system/addGame", SystemGameController.AddGame)               // 添加游戏
		// AuthApiApp.POST("/system/queryGameOne", SystemGameController.QueryGameOne)     // 查询单个游戏
		// AuthApiApp.POST("/system/editGame", SystemGameController.EditGame)             // 修改游戏
		// AuthApiApp.POST("/system/editGameStatus", SystemGameController.EditGameStatus) // 修改游戏状态

		// 游戏群管理
		// AuthApiApp.POST("/system/queryRoomList", SystemRoomController.QueryRoomList)   // 查询群列表ok
		// AuthApiApp.POST("/system/editRoomStatus", SystemRoomController.EditRoomStatus) // 修改群状态
		// AuthApiApp.POST("/system/roomCode", SystemRoomController.RoomCode)             // 群枚举

		// 注单管理
		// AuthApiApp.POST("/system/queryRedOrderRecordList", RedOrderRecordController.QueryRedRecordList) // 注单列表 ok
		// AuthApiApp.POST("/system/redInfo", RedOrderRecordController.GetRedInfo)                         // 红包详情

		// 线路提成
		// AuthApiApp.POST("/system/queryLineRoyalty", SystemLineRoyaltyController.QueryLineRoyaltyList)             // 线路提成
		// AuthApiApp.POST("/system/queryLineAgencyRoyalty", SystemLineRoyaltyController.QueryLineAgencyRoyaltyList) // 代理提成

		// 公告管理
		// AuthApiApp.POST("/system/queryPostList", SystemPostController.GetAgencyPostList) // 公告列表
		// AuthApiApp.POST("/system/queryPostById", SystemPostController.QueryPostById)     // 根据id查询公告
		// AuthApiApp.POST("/system/editPostStatus", SystemPostController.EditPostStatus)   // 修改公告状态

		// 活动/广告管理
		// AuthApiApp.POST("/system/queryActiveList", SystemActivePictureController.GetAgencyActiveList) // 活动列表
		// AuthApiApp.POST("/system/queryActiveById", SystemActivePictureController.QueryActiveById)     // 根据id查询活动
		// AuthApiApp.POST("/system/editActiveStatus", SystemActivePictureController.EditActiveStatus)   // 修改活动状态

		// 统计
		// AuthApiApp.POST("/system/orderStatistical", OrderStatistical.QueryOrderStatistical) // 盈利分析
	}
}
