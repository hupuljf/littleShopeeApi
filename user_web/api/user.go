package api

import (
	"api/user_web/global/response"
	"api/user_web/proto"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页")
	ip := "127.0.0.1"
	port := 50051
	//拨号建立连接
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户GRPC服务失败", "msg", err.Error())
	}
	//建立一个客户端访问grpc服务
	userSvcClient := proto.NewUserClient(userConn)
	rsp, err := userSvcClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 10,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList]访问用户grpc服务失败", "msg", err.Error())
		//访问grpc服务的响应的状态码 err有code信息
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, val := range rsp.Data {
		data := response.UserResponse{
			Id:       val.Id,
			NickName: val.NickName,
			Gender:   val.Gender,
			Mobile:   val.Mobile,
			Birthday: time.Time(time.Unix(int64(val.BirthDay), 0)).Format("2006-09-09"),
		}
		//data := make(map[string]interface{})
		//data["id"] = val.Id
		//data["name"] = val.NickName
		//data["birthday"] = val.BirthDay
		//data["gender"] = val.Gender
		//data["mobile"] = val.Mobile
		result = append(result, data)
	}
	ctx.JSON(http.StatusOK, result)

	//zap.S().Infof()

}
