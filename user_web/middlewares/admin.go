package middlewares

import (
	"api/user_web/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort() //在中间件处拦截
			return
		}
		ctx.Next()
	}

}
