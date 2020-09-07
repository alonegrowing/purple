package main

import (
	"purple/stone/redis"
	log "purple/stone/logging"
)

var (
	r *redis.Redis
)

func init() {
	var err error
	r, err = redis.NewRedis(&redis.RedisConfig{
		ServerName: "test",
		Addr: "localhost:6379",
		MaxIdle: 100,
		MaxActive: 100,
		IdleTimeout: 0,
		ConnectTimeout: 200,
		ReadTimeout: 100,//ms
		WriteTimeout: 100,//ms
		Database: 0,
	})
	if err != nil {
		log.Fatalf("init: %s\n", err)
	}
}

func main() {
	reply, err := r.Do("config", "get", "*")
	if err != nil {
		log.Fatalf("Do: %s\n", err)
	}
	m, _ := redis.StringMap(reply, err)
	for k,v := range m {
		log.Infof("%s:%s", k,v)
	}
}
