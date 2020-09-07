package main

import (
	"context"
	"net"
	log "purple/stone/logging"
	"purple/gen-go/purple"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func getHomePage(ctx context.Context, in *purple.ParamHomePage) (*purple.ResHomePage, error) {
	return &purple.ResHomePage{Id: in.Id}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	purple.RegisterPurpleService(s, &purple.PurpleService{GetHomePage: getHomePage})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
