package initialize

import (
	"api/user_web/global"
	"api/user_web/proto"
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

func InitUserClient() {
	//拨号建立连接
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerFromConfig.ConsulInfo.Host, global.ServerFromConfig.ConsulInfo.Port,
			global.ServerFromConfig.UserInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //多个实例 负载均衡
	)
	if err != nil {
		zap.S().Error("[GetUserList]连接用户GRPC服务失败", "msg", err.Error())
	}
	//defer conn.Close()
	//原始的配置方法取到的ip和端口号
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userIP, userPort), grpc.WithInsecure())
	userClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userClient
}

func InitUserClientV0() {
	//从consul中获取微服务 服务发现
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerFromConfig.ConsulInfo.Host, global.ServerFromConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	userIP := "" //计划ip和port由服务发现获取
	userPort := 0
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerFromConfig.UserInfo.Name)) //name
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userIP = value.Address
		userPort = value.Port
		//fmt.Println(key)
	}
	zap.S().Infof("用户服务发现ip与端口是%s:%d", userIP, userPort)
	//ip := global.ServerFromConfig.UserInfo.Host
	//port := global.ServerFromConfig.UserInfo.Port
	//拨号建立连接
	userConn, err := grpc.Dial(
		"consul://127.0.0.1:8500/user_srv?wait=14s&tag=grpc",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //多个实例 负载均衡
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer conn.Close()
	//原始的配置方法取到的ip和端口号
	//userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userIP, userPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]连接用户GRPC服务失败", "msg", err.Error())
	}
	global.UserSrvClient = proto.NewUserClient(userConn)
}
