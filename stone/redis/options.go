package redis

import (
	"fmt"
)

type RedisConfig struct {
	// This option is used to stat upload
	ServerName string `json:"server_name"`

	// Redis server host and port "localhost:6379"
	Addr string `json:"addr"`

	// Specifies the password to use when connecting to the Redis server.
	Password string `json:"password"`

	// Maximum number of idle connections in the pool.
	MaxIdle int `json:"max_idle"`

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int `json:"max_active"`

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout int `json:"idle_timeout"`

	// Specifies the timeout for connecting to the Redis server.
	ConnectTimeout int `json:"connect_timeout"`
	ReadTimeout    int `json:"read_timeout"`

	// WriteTimeout specifies the timeout for writing a single command.
	WriteTimeout int `json:"write_timeout"`

	// Database specifies the database to select when dialing a connection.
	Database int `json:"database"`

	SlowTime int `json:"slow_time"`

	Retry int `json:"retry"`
}

func (o *RedisConfig) init() error {
	if o.ServerName == "" {
		return fmt.Errorf("redis: ServerName not allowed empty string")
	}
	if o.Addr == "" {
		return fmt.Errorf("redis: Addr not allowed empty string")
	}
	if o.Database < 0 {
		return fmt.Errorf("redis: Database less than zero")
	}

	if o.MaxIdle < 0 {
		o.MaxIdle = 100
	}
	if o.MaxActive < 0 {
		o.MaxActive = 100
	}
	if o.IdleTimeout < 0 {
		o.IdleTimeout = 100
	}
	if o.ReadTimeout < 0 {
		o.ReadTimeout = 50
	}
	if o.WriteTimeout < 0 {
		o.WriteTimeout = 50
	}
	if o.ConnectTimeout < 0 {
		o.ConnectTimeout = 300
	}

	if o.SlowTime <= 0 {
		o.SlowTime = 100
	}

	if o.Retry < 0 {
		o.Retry = 0
	}

	return nil
}
