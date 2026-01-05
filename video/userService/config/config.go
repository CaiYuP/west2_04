package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"os"
	"west2-video/common/logs"
)

var C = InitConfig()

type Config struct {
	viper       *viper.Viper
	Sc          *ServerConfig
	Gc          *GrpcConfig
	Ec          *EtcdConfig
	MinIoConfig *MinIoConfig
	Jc          *JwtConfig
	RedisConfig *RedisConfig
	Mc          *MysqlConfig
	Dbc         *DbConfig
	UpC         *UploadConfig
}
type UploadConfig struct {
	APIURL  string `mapstructure:"api_url"`
	Timeout int    `mapstructure:"timeout"`
}
type JwtConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExp     int
	RefreshExp    int
}
type MysqlConfig struct {
	Host     string
	Port     int
	UserName string
	Password string
	Db       string
	Username string
}
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
type MinIoConfig struct {
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
}
type ServerConfig struct {
	Name string
	Addr string
}
type GrpcConfig struct {
	Addr    string
	Name    string
	Version string
	Weight  int64
}
type EtcdConfig struct {
	Addrs []string
}

func (c *Config) ReloadConfig() {
	c.ReadServerConfig()
	c.InitZapLog()
	c.ReadEtcdConfig()
	c.ReadMinIoConfig()
	c.InitMysqlConfig()
	c.InitJwtConfig()
	c.InitDbConfig()
	c.ReConnMysql()
	c.ReConnRedis()
	c.ReadGrpcConfig()
	c.ReadUploadConfig()
}

type DbConfig struct {
	Master          MysqlConfig
	Slave           []MysqlConfig
	Separation      bool
	MaxIdleConns    int // 最大空闲连接数
	MaxOpenConns    int // 最大打开连接数
	ConnMaxLifetime int // 连接最大生命周
}

func InitConfig() *Config {
	v := viper.New()
	conf := &Config{viper: v}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	conf.viper.AddConfigPath("/userService/config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	conf.ReloadConfig()
	return conf
}
func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{
		Name: c.viper.GetString("server.name"),
		Addr: c.viper.GetString("server.addr"),
	}
	c.Sc = sc
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
func (c *Config) InitRedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"), // no password set
		DB:       c.viper.GetInt("redis.db"),          // use default DB
	}
}
func (c *Config) InitMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.username"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetInt("mysql.port"),
		Db:       c.viper.GetString("mysql.db"),
	}
	c.Mc = mc
}
func (c *Config) ReadEtcdConfig() *EtcdConfig {
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		logs.LG.Error(err.Error())
		log.Fatalln(err)
	}
	ec := &EtcdConfig{
		Addrs: addrs,
	}
	c.Ec = ec
	return ec
}

func (c *Config) InitJwtConfig() {
	jc := &JwtConfig{
		AccessSecret:  c.viper.GetString("jwt.accessSecret"),
		RefreshSecret: c.viper.GetString("jwt.refreshSecret"),
		AccessExp:     c.viper.GetInt("jwt.accessExp"),
		RefreshExp:    c.viper.GetInt("jwt.refreshExp"),
	}
	c.Jc = jc
}

func (c *Config) InitDbConfig() {
	mc := &DbConfig{}
	mc.Separation = c.viper.GetBool("db.separation")
	var slaves []MysqlConfig
	err := c.viper.UnmarshalKey("db.slave", &slaves)
	if err != nil {
		panic(err)
	}
	master := MysqlConfig{
		Username: c.viper.GetString("db.master.username"),
		Password: c.viper.GetString("db.master.password"),
		Host:     c.viper.GetString("db.master.host"),
		Port:     c.viper.GetInt("db.master.port"),
		Db:       c.viper.GetString("db.master.db"),
	}
	mc.Master = master
	mc.Slave = slaves
	mc.MaxOpenConns = c.viper.GetInt("db.master.maxOpenConns")
	mc.MaxIdleConns = c.viper.GetInt("db.master.maxIdleConns")
	mc.ConnMaxLifetime = c.viper.GetInt("db.master.connMaxLifetime")
	c.Dbc = mc
}
func (c *Config) ReadMinIoConfig() {
	c.viper.UnmarshalKey("minio", &c.MinIoConfig)
}
func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{
		Addr:    c.viper.GetString("grpc.addr"),
		Name:    c.viper.GetString("grpc.name"),
		Version: c.viper.GetString("grpc.version"),
		Weight:  c.viper.GetInt64("grpc.weight"),
	}
	c.Gc = gc
}

func (c *Config) ReadUploadConfig() {
	c.viper.UnmarshalKey("upload", &c.UpC)
}
