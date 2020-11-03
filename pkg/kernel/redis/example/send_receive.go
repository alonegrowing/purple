package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"github.com/alonegrowing/purple/pkg/kernel/redis"
	"time"
)

var (
	r *redis.Redis
)

func init() {
	var err error
	r, err = redis.NewRedis(&redis.RedisConfig{
		ServerName:     "test",
		Addr:           "localhost:6379",
		MaxIdle:        100,
		MaxActive:      100,
		IdleTimeout:    1000,
		ConnectTimeout: 200,
		ReadTimeout:    2000, //ms
		WriteTimeout:   100,  //ms
		Database:       0,
	})
	if err != nil {
		log.Fatalf("init: %s\n", err)
	}
}

type Struct struct {
	i int
}

func (s *Struct) String() string {
	return fmt.Sprintf("this is %d", s.i)
}

func main() {
	closes := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		close(closes)
	}()

	go func() {
		for i := 0; ; i++ {
			if err := r.Send("queue_name", "this is value", i, &Struct{i: i}); err != nil {
				log.Fatalf("err:%s", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()

	bufSize := 10
	// "queue_name" is redis queue, "closes" is used to close queue,
	// bufSize is the size of channal buffer size
	for b := range r.Receive("queue_name", closes, bufSize) {
		b := b
		go func() {
			log.Printf("receive=%s\n", b)
		}()
	}
}
