package global

import (
	"api/goods_web/config"
	"api/goods_web/proto"
	ut "github.com/go-playground/universal-translator"
)

//在这里定义全局的常量
var (
	ServerFromConfig = &config.ServerConfig{}
	Trans            ut.Translator
	GoodsSrvClient   proto.GoodsClient
)
