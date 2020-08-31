package config

import (
	"purple/stone/sql"
	"purple/stone/tomlconfig"
)

var (
	logRotate  = "hour" // hour | day
	ServiceConfig Config
)


func init() {
	var conf  = "./conf/prod/config.toml"
	tomlconfig.ParseTomlConfig(conf, &ServiceConfig)

	DefaultLoggerConfig.InitLoggerConfig(ServiceConfig.Logger)
}

type LoggerConfig struct {
	Rotate          string `toml:"rotate"`
	Level           string `toml:"level"`
	LogPath         string `toml:"logpath"`
}

type Config struct {
	ServiceName string 					`toml:"service_name"`
	Logger 		LoggerConfig  			`toml:"log"`
	Database 	[]sql.SQLGroupConfig 	`toml:"database"`
}

