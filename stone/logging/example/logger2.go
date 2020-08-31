package main

import (
	"git.inke.cn/BackendPlatform/golang/logging"
)

var (
	log *logging.Logger
)

func init() {
    log = logging.NewLogger(&logging.Options{
		TimesFormat: logging.TIMESECOND,
	})
}

func main() {
	log.Infof("This message will print into stdout")
}

