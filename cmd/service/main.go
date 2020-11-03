package main

import (
	"fmt"
	"net"
	"github.com/alonegrowing/purple/gen-go/purple"
	"github.com/alonegrowing/purple/pkg/basic/util"
	"github.com/alonegrowing/purple/pkg/config"
	"github.com/alonegrowing/purple/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.Service.RPCPort))
	util.PanicIfError(err)

	srv := grpc.NewServer()
	purple.RegisterPurpleService(srv, &purple.PurpleService{
		GetHomePage: service.GetHomePage,
		GetMember:   service.GetMember,
	})
	err = srv.Serve(listener)
	util.PanicIfError(err)
}
