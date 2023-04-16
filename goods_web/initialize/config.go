package initialize

import (
	"api/goods_web/global"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	debug := GetEnvInfo("CONFIG")
	configFilePrefix := "config"
	fmt.Println(debug)
	configFileName := fmt.Sprintf("goods_web/%s-pro.yaml", configFilePrefix)
	if debug == "test" {
		configFileName = fmt.Sprintf("goods_web/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//ServerFromConfig 得是一个全局变量
	//global.ServerFromConfig = &config.ServerConfig{}
	if err := v.Unmarshal(global.ServerFromConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerFromConfig)

	//viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("config file changed: %s ", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerFromConfig)
		zap.S().Infof("现在的配置信息：%v", global.ServerFromConfig)
	})

	//time.Sleep(time.Second * 300)

}

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}
