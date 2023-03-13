package initialize

import (
	"api/user_web/middlewares"
	"api/user_web/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	//返回传给服务器的总路由 /v1/user/...
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup) // /u/v1/user/....
	router.InitBaseRouter(ApiGroup) // /u/v1/base/....
	return Router
}
