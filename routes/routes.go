package routes

import (
	"baseGo/controller"
	"baseGo/middleware"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CROSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	return r
}
