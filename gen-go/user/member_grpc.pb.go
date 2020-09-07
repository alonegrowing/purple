// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// MemberClient is the client API for Member service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MemberClient interface {
	GetMember(ctx context.Context, in *GetMemberParam, opts ...grpc.CallOption) (*MemberResponse, error)
}

type memberClient struct {
	cc grpc.ClientConnInterface
}

func NewMemberClient(cc grpc.ClientConnInterface) MemberClient {
	return &memberClient{cc}
}

var memberGetMemberStreamDesc = &grpc.StreamDesc{
	StreamName: "GetMember",
}

func (c *memberClient) GetMember(ctx context.Context, in *GetMemberParam, opts ...grpc.CallOption) (*MemberResponse, error) {
	out := new(MemberResponse)
	err := c.cc.Invoke(ctx, "/user.Member/GetMember", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MemberService is the service API for Member service.
// Fields should be assigned to their respective handler implementations only before
// RegisterMemberService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type MemberService struct {
	GetMember func(context.Context, *GetMemberParam) (*MemberResponse, error)
}

func (s *MemberService) getMember(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMemberParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/user.Member/GetMember",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetMember(ctx, req.(*GetMemberParam))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterMemberService registers a service implementation with a gRPC server.
func RegisterMemberService(s grpc.ServiceRegistrar, srv *MemberService) {
	srvCopy := *srv
	if srvCopy.GetMember == nil {
		srvCopy.GetMember = func(context.Context, *GetMemberParam) (*MemberResponse, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetMember not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "user.Member",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetMember",
				Handler:    srvCopy.getMember,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "member.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewMemberService creates a new MemberService containing the
// implemented methods of the Member service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewMemberService(s interface{}) *MemberService {
	ns := &MemberService{}
	if h, ok := s.(interface {
		GetMember(context.Context, *GetMemberParam) (*MemberResponse, error)
	}); ok {
		ns.GetMember = h.GetMember
	}
	return ns
}

// UnstableMemberService is the service API for Member service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableMemberService interface {
	GetMember(context.Context, *GetMemberParam) (*MemberResponse, error)
}
