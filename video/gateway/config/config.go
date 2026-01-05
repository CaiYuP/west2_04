package config

import (
	"log"
	"os"
	"west2-video/common/logs"

	"github.com/spf13/viper"
)

var C = InitConfig()

type Config struct {
	viper    *viper.Viper
	Server   ServerConfig   `mapstructure:"server"`
	Services ServicesConfig `mapstructure:"services"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type ServicesConfig struct {
	UserService        string `mapstructure:"user_service"`
	VideoService       string `mapstructure:"video_service"`
	InteractionService string `mapstructure:"interaction_service"`
	SocialService      string `mapstructure:"social_service"`
}

type JWTConfig struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpiresIn  int64
	RefreshExpiresIn int64
}

var globalConfig *Config

// InitConfig 初始化配置
func InitConfig() *Config {
	v := viper.New()

	workDir, _ := os.Getwd()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(workDir + "/config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	cfg.viper = v
	globalConfig = &cfg
	cfg.InitZapLog()
	log.Printf("配置文件加载成功: %s", v.ConfigFileUsed())
	return &cfg
}
func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}
