package global

import (
	"api/user_web/config"
	"api/user_web/proto"
	ut "github.com/go-playground/universal-translator"
)

//在这里定义全局的常量
var (
	ServerFromConfig = &config.ServerConfig{}
	Trans            ut.Translator
	UserSrvClient    proto.UserClient
)
