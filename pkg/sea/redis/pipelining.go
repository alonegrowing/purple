package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"sync"
)

type Pipelining struct {
	conn redis.Conn
	*sync.Mutex
	isClose bool
}

func (r *Redis) NewPipelining() (*Pipelining, error) {
	p := &Pipelining{}
	client := r.pool.Get()
	err := client.Err()
	if err != nil {
		return nil, err
	}
	p.conn = client
	p.Mutex = &sync.Mutex{}
	return p, nil
}

func (p *Pipelining) Send(cmd string, args ...interface{}) error {
	p.Lock()
	defer p.Unlock()
	if p.isClose {
		return errors.New("Pipelining closed")
	}
	return p.conn.Send(cmd, args...)
}

func (p *Pipelining) Flush() error {
	p.Lock()
	defer p.Unlock()
	if p.isClose {
		return errors.New("Pipelining closed")
	}
	return p.conn.Flush()
}

func (p *Pipelining) Receive() (reply interface{}, err error) {
	p.Lock()
	defer p.Unlock()
	if p.isClose {
		return nil, errors.New("Pipelining closed")
	}
	return p.conn.Receive()
}

func (p *Pipelining) Close() error {
	p.Lock()
	defer p.Unlock()
	p.isClose = true
	return p.conn.Close()
}
