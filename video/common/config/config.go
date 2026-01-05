package config

import (
	"github.com/spf13/viper"
	"log"
)

var C = InitConfig()

type Config struct {
	viper *viper.Viper
	UpC   *UploadConfig
}

func (c *Config) ReloadConfig() {
	c.ReadUploadConfig()
}

func (c *Config) ReadUploadConfig() {
	c.viper.UnmarshalKey("upload", &c.UpC)
}

type UploadConfig struct {
	APIURL  string `mapstructure:"api_url"`
	Timeout int    `mapstructure:"timeout"`
}

func InitConfig() *Config {
	v := viper.New()
	conf := &Config{viper: v}
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath("E:/west2/04/west2_04/video/common/config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	conf.ReloadConfig()
	return conf
}
