package main

import (
	log "purple/stone/logging"
	"purple/stone/redis"
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
	reply, err := r.Do("TIME")
	if err != nil {
		log.Fatalf("Do: %s\n", err)
	}
	ss, _ := redis.Strings(reply, err)
	log.Infof("stirng1:%s string2:%s", ss[0], ss[1])
}
