package middlewares

import (
	"api/goods_web/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentgoods := claims.(*models.CustomClaims)
		fmt.Println(currentgoods)
		if currentgoods.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort() //在中间件处拦截
			return
		}
		ctx.Next()
	}

}
