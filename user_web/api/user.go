package api

import (
	"api/user_web/forms"
	"api/user_web/global"
	"api/user_web/global/response"
	"api/user_web/middlewares"
	"api/user_web/models"
	"api/user_web/proto"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
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
					"error": e.Code(),
				})
			}
			return
		}
	}
}

func GetUserList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims") //interface{}类型
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户ID：%d", currentUser.ID)

	//pn和psize可以代进query里面 由ctx取到
	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("psize", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)
	//建立一个客户端访问grpc服务
	//userSvcClient := proto.NewUserClient(userConn)
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
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

func PassWordLogin(ctx *gin.Context) {
	var loginForm forms.PassWordLoginForm
	if err := ctx.ShouldBind(&loginForm); err != nil {
		errs, ok := err.(validator.ValidationErrors) //将error强制转换为validationerror类型
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": removeTopStruct(errs.Translate(global.Trans)), //让error可以发中文
		})
		return
	}
	//连接grpc服务前处理验证码
	if !store.Verify(loginForm.CaptchaId, loginForm.Captcha, true) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	userSvcClient := global.UserSrvClient
	if rsp, err := userSvcClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: loginForm.Mobile,
	}); err != nil { //没找到的情况下 报的error
		fmt.Println(err.Error())
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "该用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "出错了",
				})

			}
		}
	} else { //查询到用户之后处理密码校验
		if passRsp, passErr := userSvcClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          loginForm.PassWord,
			EncryptedPassword: rsp.PassWord,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "登陆失败",
			})
		} else {
			if passRsp.Success {
				//登录完成后生成jwt的token
				jsonWebToken := middlewares.NewJWT()
				claims := models.CustomClaims{ //生成token的一些参数
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               //签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, //单位是秒 30天后过期
						Issuer:    "runzheng",
					},
				}
				token, err := jsonWebToken.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, map[string]string{
						"msg": "生成token失败",
					})
					return
				}
				//zap.S().Error(token)
				ctx.JSON(http.StatusOK, gin.H{
					"msg": "登陆成功",
					//"id":         rsp.Id,
					//"nickname":   rsp.NickName,
					"data":       rsp,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //ms

				})

			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"msg": "密码错误",
				})
			}

		}
	}

	return

}

func Register(ctx *gin.Context) {
	var registerForm forms.RegisterForm
	if err := ctx.ShouldBind(&registerForm); err != nil {
		errs, ok := err.(validator.ValidationErrors) //将error强制转换为validationerror类型
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": removeTopStruct(errs.Translate(global.Trans)), //让error可以发中文
		})
		return
	}
	//验证码处理逻辑：

	//grpc调用
	userSvcClient := global.UserSrvClient
	rsp, err := userSvcClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   registerForm.Mobile,
		PassWord: registerForm.PassWord,
		NickName: "bobby" + registerForm.Mobile[6:],
	})
	if err != nil {
		zap.S().Errorf("[Register]【新建用户失败】失败: %s", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		//HandleGrpcErrorToHttp(err, ctx)
		return
	}
	jsonWebToken := middlewares.NewJWT()
	claims := models.CustomClaims{ //生成token的一些参数
		ID:          uint(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: uint(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //单位是秒 30天后过期
			Issuer:    "runzheng",
		},
	}
	token, err := jsonWebToken.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg": "生成token失败",
		})
		return
	}
	//zap.S().Error(token)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "注册成功",
		//"id":         rsp.Id,
		//"nickname":   rsp.NickName,
		"data":       rsp,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000, //ms

	})
	return

}
