package main

import (
	"api/user_web/global"
	"api/user_web/initialize"
	validator2 "api/user_web/validator"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func main() {
	//port := 8021
	initialize.InitLogger()
	initialize.InitConfig()
	//初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Errorf("初始化错误翻译器错误 %s", err.Error())
		panic(err)
	}
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile_validate", validator2.ValidateMobile)
		_ = v.RegisterTranslation("mobile_validate", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile_validate", "{0} 非法的手机号码!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile_validate", fe.Field())
			return t
		})
	}
	//1.初始化router
	r := initialize.Routers()
	/*
		1.S()获取一个全局的sugar 但是自己一点都不做 咱们可以自己设置logger
	*/
	//logger, _ := zap.NewProduction()
	//zap.ReplaceGlobals(logger)
	zap.S().Infof("启动服务器，端口：%d", global.ServerFromConfig.Port)
	if err := r.Run(fmt.Sprintf(":%d", global.ServerFromConfig.Port)); err != nil {
		//err.Error()才有值
		zap.S().Panicf("启动失败%s", err.Error())
	}
}
