package main

import (
	"log"
	"github.com/alonegrowing/purple/pkg/kernel/redis"
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
		IdleTimeout:    0,
		ConnectTimeout: 200,
		ReadTimeout:    100, //ms
		WriteTimeout:   100, //ms
		Database:       0,
	})
	if err != nil {
		log.Fatalf("init: %s\n", err)
	}
}

func main() {
	_, err := r.Do("BLPOP", "NotExistKey", 1)
	if err != redis.ErrTimeout {
		log.Fatalf("Do: %s\n", err)
	}
}
