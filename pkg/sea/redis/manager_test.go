package redis

import (
	"testing"

	"log"
	_ "time"
)

var (
	m *Manager
)

func init() {
	var err error

	c := []RedisConfig{}
	c = append(c, RedisConfig{
		Addr:       "localhost:6379",
		ServerName: "haha1",
	})
	c = append(c, RedisConfig{
		Addr:       "localhost:6379",
		ServerName: "haha2",
	})

	m, err = NewManager(c)

	if err != nil {
		log.Fatalf("err:%s", err)
	}
}

func TestManager(t *testing.T) {
	r1 := m.Get("haha1")
	r1.Set("set1", 1)

	//	r2 := m.Get("haha2")
	//	r2.Set("set1", 2)
}
