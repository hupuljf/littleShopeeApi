package global

import (
	"api/user_web/config"
	ut "github.com/go-playground/universal-translator"
)

//在这里定义全局的常量
var (
	ServerFromConfig = &config.ServerConfig{}
	Trans            ut.Translator
)
