package routes

import (
	"baseGo/controller"
	"baseGo/middleware"

	"github.com/gin-gonic/gin"
)

var(
	SystemAgencyController        = new(controller.SystemAgencyController)
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	// r.Use(middleware.CROSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	AuthApiApp := r.Group("/api"){

		// 业主超管管理
		AuthApiApp.POST("/system/queryAgencyAdminList", SystemAgencyController.QueryAgencyAdminList) // 查询超管列表ok
		AuthApiApp.POST("/system/addAgencyAdmin", SystemAgencyController.AddAgencyAdmin)             // 添加超管
		AuthApiApp.POST("/system/queryAgencyAdminOne", SystemAgencyController.QueryAgencyAdminOne)   // 查询单个超管
		AuthApiApp.POST("/system/editAgencyAdmin", SystemAgencyController.EditAgencyAdmin)           // 修改超管


	}
	return r
}
