package router

import (
	"api/goods_web/api/goods"
	"api/goods_web/middlewares"
	"github.com/gin-gonic/gin"
)

func InitgoodsRouter(Router *gin.RouterGroup) {
	goodsRouter := Router.Group("goods")
	//增加大括号 增加可读性
	{
		goodsRouter.GET("list", goods.GetGoodsList)
		goodsRouter.POST("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.NewGoods)
		goodsRouter.GET("health")

	}

}
