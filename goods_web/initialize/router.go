package initialize

import (
	"api/goods_web/middlewares"
	"api/goods_web/router"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	//返回传给服务器的总路由 /v1/goods/...
	Router := gin.Default()
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/g/v1")
	router.InitgoodsRouter(ApiGroup) // /g/v1/goods/....
	return Router
}
