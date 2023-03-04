package router

import (
	"api/user_web/api"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	//增加大括号 增加可读性
	{
		UserRouter.GET("list", api.GetUserList)

	}

}
