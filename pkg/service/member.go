package service

import (
	"context"
	"github.com/alonegrowing/purple/gen-go/purple"
)

func GetMember(ctx context.Context, in *purple.GetMemberParam) (*purple.MemberResponse, error) {
	return &purple.MemberResponse{Id: in.Id}, nil
}
