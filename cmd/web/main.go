package main

import (
	"fmt"
	"github.com/alonegrowing/purple/pkg/config"
	"github.com/alonegrowing/purple/pkg/web"
)

func main() {
	route := web.NewRouter()
	_ = route.Run(fmt.Sprintf(":%d", config.ServiceConfig.Service.WEBPort))
}
