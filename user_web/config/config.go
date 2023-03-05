package config

type UserConfig struct {
	Host string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	ServiceName string     `mapstructure:"name"`
	Port        int        `mapstructure:"port"`
	UserInfo    UserConfig `mapstructure:"user_srv"`
}
