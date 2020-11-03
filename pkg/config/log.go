package config

import (
	log "github.com/alonegrowing/purple/pkg/kernel/logging"
)

type LoggerConfig struct {
	Rotate  string `toml:"rotate"`
	Level   string `toml:"level"`
	LogPath string `toml:"logpath"`
}

func InitLoggerConfig(conf LoggerConfig) {
	if conf.Rotate == logRotate {
		log.SetRotateByHour()
	} else {
		log.SetRotateByDay()
	}
	log.SetOutputPath(conf.LogPath)
	log.SetLevelByString(conf.Level)
}
