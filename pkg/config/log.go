package config

import log "purple/stone/logging"


type LoggerConfiger interface{
	InitLoggerConfig(conf LoggerConfig)
}

type LoggerConfigerImpl struct {}

var DefaultLoggerConfig LoggerConfiger

func init() {
	DefaultLoggerConfig = NewLoggerConfigerImpl()
}

func NewLoggerConfigerImpl() *LoggerConfigerImpl {
	return &LoggerConfigerImpl{
	}
}

func (r *LoggerConfigerImpl) InitLoggerConfig(conf LoggerConfig) {
	if conf.Rotate == logRotate {
		log.SetRotateByHour()
	} else {
		log.SetRotateByDay()
	}
	log.SetOutputPath(conf.LogPath)
	log.SetLevelByString(conf.Level)
}