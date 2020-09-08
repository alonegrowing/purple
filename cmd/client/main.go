package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"purple/gen-go/purple"
	"purple/pkg/config"
	"purple/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", config.ServiceConfig.Service.RPCPort), grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	utils.PanicIfError(err)

	c := purple.NewPurpleClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r1, err := c.GetHomePage(ctx, &purple.HomePageParam{Id: 1226})
	utils.PanicIfError(err)
	log.Printf("Greeting: %d", r1.GetId())

	r2, err := c.GetMember(ctx, &purple.GetMemberParam{Id: 111})
	utils.PanicIfError(err)
	log.Printf("Greeting: %d", r2.GetId())
}
