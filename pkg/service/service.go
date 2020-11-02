package service

import (
	"context"
	"github.com/alonegrowing/purple/gen-go/purple"
)

func GetHomePage(ctx context.Context, in *purple.HomePageParam) (*purple.HomePageResponse, error) {
	return &purple.HomePageResponse{Id: in.Id}, nil
}
