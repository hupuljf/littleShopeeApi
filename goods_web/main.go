package main

import (
	"api/goods_web/global"
	"api/goods_web/initialize"
	"api/goods_web/utils"
	"api/goods_web/utils/register/consul"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	//port := 8021
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitgoodsClient() //连接后台微服务的客户端
	//初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Errorf("初始化错误翻译器错误 %s", err.Error())
		panic(err)
	}
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
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
	//zap.S().Infof("启动服务器，端口：%d", global.ServerFromConfig.Port)

	//如果env是pro 就用随机端口(修改初始化配置中的port)：
	if initialize.GetEnvInfo("CONFIG") != "test" {
		global.ServerFromConfig.Port, _ = utils.GetFreePort()
	}
	//注册服务
	register_client := consul.NewRegistry(global.ServerFromConfig.ConsulInfo.Host, global.ServerFromConfig.ConsulInfo.Port)

	uuid := fmt.Sprintf("%s", uuid.NewV4())
	register_client.Register("192.168.2.9", 8021, "goods_web", make([]string, 0), uuid)

	zap.S().Infof("启动服务器，端口：%d", global.ServerFromConfig.Port)
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", global.ServerFromConfig.Port)); err != nil {
			//err.Error()才有值
			zap.S().Panicf("启动失败%s", err.Error())
		}
	}() //得有go func放进协程里面 后面的程序才回去执行

	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //终止或者异常退出的信号量打给quit 会去监听 一旦监听到了 打给quit
	<-quit                                               //阻塞
	if register_client.DeRegister(uuid) != nil {
		zap.S().Error("注销失败")
	} else {
		zap.S().Error("注销成功")
	}

}
