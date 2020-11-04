package resource

import (
	"errors"
	"github.com/alonegrowing/purple/pkg/sea/redis"
)

var (
	REDIS_CLIENT_NOT_INIT = errors.New("redis client not init")
)

var defaultRedis map[string]*redis.Redis

func NewRedis(redisConfigs []redis.RedisConfig) error {
	if defaultRedis == nil {
		defaultRedis = make(map[string]*redis.Redis)
	}
	for _, conf := range redisConfigs {
		client, err := redis.NewRedis(&conf)
		if err != nil || client == nil {
			continue
		}
		defaultRedis[conf.ServerName] = client
	}
	return nil
}

func GetRedis(service string) (*redis.Redis, error) {
	if client, ok := defaultRedis[service]; ok {
		return client, nil
	}
	return nil, REDIS_CLIENT_NOT_INIT
}
