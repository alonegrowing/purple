package config

import (
	"purple/stone/sql"
	"purple/stone/tomlconfig"
)

var ServiceConfig Config

func init() {
	var conf  = "./conf/prod/config.toml"
	tomlconfig.ParseTomlConfig(conf, &ServiceConfig)
}

type Config struct {
	ServiceName string `toml:"service_name"`

	Log struct {
		Level           string `toml:"level"`
		Rotate          string `toml:"rotate"`
		Accesslog       string `toml:"accesslog"`
		Businesslog     string `toml:"businesslog"`
		Serverlog       string `toml:"serverlog"`
		StatLog         string `toml:"statlog"`
		ErrorLog        string `toml:"errlog"`
		LogPath         string `toml:"logpath"`
		BalanceLogLevel string `toml:"balance_log_level"`
		GenLogLevel     string `toml:"gen_log_level"`
		Filename        string `toml:"filename"`
	} `toml:"log"`

	Database []sql.SQLGroupConfig `toml:"database"`
}
