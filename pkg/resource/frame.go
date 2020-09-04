package resource

import (
	"purple/pkg/config"
	log "purple/stone/logging"
)

func init() {
	InitFrameworkUtils(config.ServiceConfig)
}

func InitFrameworkUtils(c config.Config) {
	if c.Log.Rotate == LOG_ROTATE_HOUR {
		log.SetRotateByHour()
	} else {
		log.SetRotateByDay()
	}
	log.SetOutputPath(c.Log.LogPath)
	log.SetLevelByString(c.Log.Level)
}
