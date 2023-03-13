package config

type UserConfig struct {
	Host string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type ServerConfig struct {
	ServiceName string     `mapstructure:"name"`
	Port        int        `mapstructure:"port"`
	UserInfo    UserConfig `mapstructure:"user_srv"`
	JWTInfo     JWTConfig  `mapstructure:"jwt"`
}
