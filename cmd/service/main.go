package main

import (
	"context"
	"net"
	log "purple/stone/logging"

	//"net"

	//"google.golang.org/grpc"
	pb "purple/gen-go/user"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func GetMember(ctx context.Context, in *pb.GetMemberParam) (*pb.MemberResponse, error) {
	return &pb.MemberResponse{Id: in.Id}, nil
}

func main() {

	/*
		test := &pb.GetMemberParam{Id: 1}
		// 进行编码
		data, err := proto.Marshal(test)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}

		// 进行解码
		newTest := &pb.GetMemberParam{}
		err = proto.Unmarshal(data, newTest)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}

		fmt.Printf("id:%d", newTest.Id)

		// 测试结果
		if test.String() != newTest.String() {
			log.Fatalf("data mismatch %q != %q", test.String(), newTest.String())
		}
	*/
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMemberService(s, &pb.MemberService{GetMember: GetMember})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	/*

	 */
}
