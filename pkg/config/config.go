package config

import (
	"purple/stone/sql"
	"purple/stone/tomlconfig"
	"purple/stone/redis"
	"purple/pkg/resource"
)

var (
	logRotate     = "hour" // hour | day
	ServiceConfig Config
)

func init() {
	var conf = "./conf/prod/config.toml"

	// config parser
	tomlconfig.ParseTomlConfig(conf, &ServiceConfig)

	// logger init
	DefaultLoggerConfig.InitLoggerConfig(ServiceConfig.Logger)

	// redis init
	resource.NewRedis(ServiceConfig.Redis)

	// mysql init
	resource.NewMysqlGroup(ServiceConfig.Database)
}

type LoggerConfig struct {
	Rotate  string `toml:"rotate"`
	Level   string `toml:"level"`
	LogPath string `toml:"logpath"`
}

type Config struct {
	ServiceName string               `toml:"service_name"`
	Logger      LoggerConfig         `toml:"log"`
	Database    []sql.SQLGroupConfig `toml:"database"`
	Redis       []redis.RedisConfig  `toml:"redis"`
}
