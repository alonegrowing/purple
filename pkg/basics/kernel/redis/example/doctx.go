package main

import (
	log "purple/stone/logging"
	"purple/stone/redis"

	"context"
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
	timeout := 20 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	reply, err := r.DoCtx(ctx, "BLPOP", "NotExistKey", 0)
	if err != context.DeadlineExceeded {
		log.Fatalf("TestTimeOut err: %s %v\n", err, reply)
	}
}
