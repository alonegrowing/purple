package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"purple/gen-go/purple"
	"purple/pkg/config"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", config.ServiceConfig.Service.RPCPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := purple.NewPurpleClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetHomePage(ctx, &purple.HomePageParam{Id: 1226})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %d", r.GetId())
}
