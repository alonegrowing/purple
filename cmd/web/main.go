package main

import (
	"fmt"
	"purple/pkg/config"
	"purple/pkg/web"
)

func main() {
	route := web.NewRouter()
	_ = route.Run(fmt.Sprintf(":%d", config.ServiceConfig.Service.WEBPort))
}
