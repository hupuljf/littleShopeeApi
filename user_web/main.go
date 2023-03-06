package main

import (
	"api/user_web/global"
	"api/user_web/initialize"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	//port := 8021
	initialize.InitLogger()
	initialize.InitConfig()
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
