package main

import (
	"fmt"
	"net"
	"purple/gen-go/purple"
	"purple/pkg/config"
	"purple/pkg/service"
	log "purple/stone/logging"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServiceConfig.Service.RPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	register(server)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func register(s *grpc.Server) {
	purple.RegisterPurpleService(s, &purple.PurpleService{
		GetHomePage: service.GetHomePage,
	})
	purple.RegisterMemberService(s, &purple.MemberService{
		GetMember: service.GetMember,
	})
}
