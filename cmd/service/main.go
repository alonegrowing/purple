package main

import (
	"fmt"
	"net"
	"purple/gen-go/purple"
	"purple/pkg/config"
	"purple/pkg/service"
	"google.golang.org/grpc"
	"purple/pkg/utils"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.Service.RPCPort))
	utils.PanicIfError(err)

	srv := grpc.NewServer()
	purple.RegisterPurpleService(srv, &purple.PurpleService{
		GetHomePage: service.GetHomePage,
		GetMember: service.GetMember,
	})
	err = srv.Serve(listener)
	utils.PanicIfError(err)
}
