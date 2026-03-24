package config

import (
	"github.com/spf13/viper"
	"log"
)

//var BC = initBootConfit()

type BootConfig struct {
	viper       *viper.Viper
	NacosConfig *NacosConfig
}
type NacosConfig struct {
	Namespace   string
	Group       string
	IpAddr      string
	Port        int
	ContextPath string
	Scheme      string
}

func (bc *BootConfig) InitNacosConfig() {
	nc := &NacosConfig{}
	bc.viper.UnmarshalKey("nacos", nc)
	bc.NacosConfig = nc
}
func InitBootConfit() *BootConfig {
	bc := &BootConfig{}
	bc.viper = viper.New()
	bc.viper.SetConfigName("bootstrap")
	bc.viper.SetConfigType("yaml")
	bc.viper.AddConfigPath("E:/west2/04/west2_04/video/socialService/config")
	err := bc.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	bc.InitNacosConfig()
	return bc
}
