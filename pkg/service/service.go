package service

import (
	"context"
	"purple/gen-go/purple"
)

func GetHomePage(ctx context.Context, in *purple.HomePageParam) (*purple.HomePageResponse, error) {
	return &purple.HomePageResponse{Id: in.Id}, nil
}
