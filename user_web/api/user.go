package api

import (
	"api/user_web/global"
	"api/user_web/global/response"
	"api/user_web/proto"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var trans ut.Translator

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
	ip := global.ServerFromConfig.UserInfo.Host
	port := global.ServerFromConfig.UserInfo.Port
	//拨号建立连接
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户GRPC服务失败", "msg", err.Error())
	}
	//pn和psize可以代进query里面 由ctx取到
	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)
	//建立一个客户端访问grpc服务
	userSvcClient := proto.NewUserClient(userConn)
	rsp, err := userSvcClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
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
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func InitTrans(locale string) (err error) {
	//修改gin框架中的validator引擎属性, 实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
		}
		return
	}

	return
}

func PassWordLogin(ctx *gin.Context) {

}
