package initialize

import (
	"api/goods_web/global"
	"api/goods_web/proto"
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

func InitgoodsClient() {
	//拨号建立连接
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerFromConfig.ConsulInfo.Host, global.ServerFromConfig.ConsulInfo.Port,
			global.ServerFromConfig.GoodsInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //多个实例 负载均衡
	)
	if err != nil {
		zap.S().Error("[GetGoodsList]连接用户GRPC服务失败", "msg", err.Error())
	}
	//defer conn.Close()
	//原始的配置方法取到的ip和端口号
	//goodsConn, err := grpc.Dial(fmt.Sprintf("%s:%d", goodsIP, goodsPort), grpc.WithInsecure())
	goodsClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsClient
}

func InitgoodsClientV0() {
	//从consul中获取微服务 服务发现
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerFromConfig.ConsulInfo.Host, global.ServerFromConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	goodsIP := "" //计划ip和port由服务发现获取
	goodsPort := 0
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerFromConfig.GoodsInfo.Name)) //name
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		goodsIP = value.Address
		goodsPort = value.Port
		//fmt.Println(key)
	}
	zap.S().Infof("用户服务发现ip与端口是%s:%d", goodsIP, goodsPort)
	//ip := global.ServerFromConfig.goodsInfo.Host
	//port := global.ServerFromConfig.goodsInfo.Port
	//拨号建立连接
	goodsConn, err := grpc.Dial(
		"consul://127.0.0.1:8500/goods_srv?wait=14s&tag=grpc",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), //多个实例 负载均衡
	)
	if err != nil {
		log.Fatal(err)
	}
	//defer conn.Close()
	//原始的配置方法取到的ip和端口号
	//goodsConn, err := grpc.Dial(fmt.Sprintf("%s:%d", goodsIP, goodsPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetgoodsList]连接用户GRPC服务失败", "msg", err.Error())
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
}
