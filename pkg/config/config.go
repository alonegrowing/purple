package config

import (
	"github.com/alonegrowing/purple/pkg/basic/resource"
	"github.com/alonegrowing/purple/pkg/sea/redis"
	"github.com/alonegrowing/purple/pkg/sea/sql"
	"github.com/alonegrowing/purple/pkg/sea/tomlconfig"
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
	InitLoggerConfig(ServiceConfig.Logger)

	// redis init
	resource.NewRedis(ServiceConfig.Redis)

	// mysql init
	resource.NewMysqlGroup(ServiceConfig.Database)
}

type Service struct {
	WEBPort int64 `toml:"web_port"`
	RPCPort int64 `toml:"rpc_port"`
}

type Config struct {
	ServiceName string               `toml:"service_name"`
	Service     Service              `toml:"service"`
	Logger      LoggerConfig         `toml:"log"`
	Database    []sql.SQLGroupConfig `toml:"database"`
	Redis       []redis.RedisConfig  `toml:"redis"`
}
