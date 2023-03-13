package router

import (
	"api/user_web/api"
	"api/user_web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	//增加大括号 增加可读性
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList) //要经过jwt认证才能得到用户列表
		UserRouter.POST("login", api.PassWordLogin)
		UserRouter.POST("register", api.Register)
		UserRouter.GET("health")

	}

}
