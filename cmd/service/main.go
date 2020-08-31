
package main

import (
	"context"
	log "purple/stone/logging"
	"net"

	"google.golang.org/grpc"
	pb "purple/gen-go/user"
)

const (
	port = ":50051"
)

func getMember(ctx context.Context, in *pb.GetMemberParam) (*pb.MemberResponse, error) {
	return &pb.MemberResponse{Id: in.Id}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMemberService(s, &pb.MemberService{GetMember:getMember})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}