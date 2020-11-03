package redis

import (
	"context"
	"errors"
	log "github.com/alonegrowing/purple/pkg/kernel/logging"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var serverLocalPid = os.Getpid()

var logFormat = "2006/01/02 15:04:05"

type Redis struct {
	pool     *redis.Pool
	opts     *RedisConfig
	lastTime int64
}

func NewRedis(o *RedisConfig) (store *Redis, err error) {
	if err = o.init(); err != nil {
		return
	}
	opts := []redis.DialOption{}
	opts = append(opts, redis.DialConnectTimeout(time.Duration(o.ConnectTimeout)*time.Millisecond))
	opts = append(opts, redis.DialReadTimeout(time.Duration(o.ReadTimeout)*time.Millisecond))
	opts = append(opts, redis.DialWriteTimeout(time.Duration(o.WriteTimeout)*time.Millisecond))
	if len(o.Password) != 0 {
		opts = append(opts, redis.DialPassword(o.Password))
	}
	opts = append(opts, redis.DialDatabase(o.Database))
	pool := redisinit(o.Addr, o.Password, o.MaxIdle, o.IdleTimeout, o.MaxActive, opts...)
	oo := *o
	return &Redis{
		pool:     pool,
		opts:     &oo,
		lastTime: time.Now().UnixNano(),
	}, nil
}

func redisinit(server, password string, maxIdle, idleTimeout, maxActive int, options ...redis.DialOption) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		MaxActive:   maxActive,
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			protocol := "tcp"
			if strings.HasPrefix(server, "unix://") {
				server = strings.TrimLeft(server, "unix://")
				protocol = "unix"
			}
			c, err = redis.Dial(protocol, server, options...)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (r *Redis) RPop(key string) (res string, err error) {
	reply, err := r.do("RPOP", redisBytes, key)
	if err != nil {
		return "", err
	}
	return string(reply.([]byte)[:]), nil
}

func (r *Redis) LPush(name string, fields ...interface{}) error {
	keys := []interface{}{name}
	keys = append(keys, fields...)
	_, err := r.do("LPUSH", nil, keys...)
	return err
}

func (r *Redis) Send(name string, fields ...interface{}) error {
	keys := []interface{}{name}
	keys = append(keys, fields...)
	_, err := r.do("RPUSH", nil, keys...)
	return err
}

func (r *Redis) Receive(name string, closech chan struct{}, bufferSize int) chan []byte {
	ch := make(chan []byte, bufferSize)
	go func() {
		defer close(ch)
		for {
			select {
			case <-closech:
				return
			default:
				data, err := r.do("BLPOP", nil, name, 1)
				if err == nil {
					if data != nil {
						ms, err := redis.ByteSlices(data, nil)
						if err != nil {
							log.Errorf("convert redis response error %v", err)
						} else {
							ch <- ms[1]
						}
					}
				} else if err != ErrTimeout {
					log.Errorf("BRPOP error %s", err)
				}
			}
		}
	}()
	return ch
}

func (r *Redis) Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	return r.do(cmd, nil, args...)
}

func (r *Redis) DoCtx(ctx context.Context, cmd string, args ...interface{}) (reply interface{}, err error) {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		reply, err = r.do(cmd, nil, args...)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-ch:
		return
	}
}

func (r *Redis) Set(key, value interface{}) (ret bool, err error) {
	var reply interface{}
	reply, err = r.do("SET", redisString, key, value)
	if err != nil {
		return
	}
	rsp := reply.(string)

	if rsp == "OK" {
		ret = true
	}

	return
}

func (r *Redis) SetExSecond(key, value interface{}, dur int) (ret string, err error) {
	var reply interface{}
	reply, err = r.do("SET", redisString, key, value, "EX", dur)
	if err != nil {
		return
	}
	ret = reply.(string)
	return
}

func (r *Redis) Get(key string) (ret []byte, err error) {
	var reply interface{}
	reply, err = r.do("GET", redisBytes, key)
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			var tmp []byte
			ret = tmp
		}
		return
	}
	ret = reply.([]byte)
	return
}

func (r *Redis) GetInt(key string) (ret int, err error) {
	var reply interface{}
	reply, err = r.do("GET", redisInt, key)
	if err != nil {
		return
	}
	ret = reply.(int)
	return
}

func (r *Redis) MGet(keys ...interface{}) (ret [][]byte, err error) {
	var reply interface{}
	reply, err = r.do("MGET", redisByteSlices, keys...)
	if err != nil {
		return
	}
	ret = reply.([][]byte)
	return
}

func (r *Redis) MSet(keys ...interface{}) (ret string, err error) {
	var reply interface{}
	reply, err = r.do("MSET", redisString, keys...)
	if err != nil {
		return
	}
	ret = reply.(string)
	return
}

func (r *Redis) Del(args ...interface{}) (count int, err error) {
	var reply interface{}
	reply, err = r.do("Del", redisInt, args...)
	if err != nil {
		return
	}
	count = reply.(int)
	return
}

func (r *Redis) Exists(key string) (res bool, err error) {
	var reply interface{}
	reply, err = r.do("Exists", redisBool, key)
	if err != nil {
		return
	}
	res = reply.(bool)
	return
}

func (r *Redis) Expire(key string, expire time.Duration) error {
	_, err := r.do("EXPIRE", nil, key, expire.Seconds())
	if err != nil {
		return err
	}
	return nil
}

/*
*	hash
 */
func (r *Redis) HDel(key interface{}, fields ...interface{}) (res int, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, fields...)

	reply, err = r.do("HDEL", redisInt, keys...)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) HSet(key, fieldk string, fieldv interface{}) (res int, err error) {
	var reply interface{}
	reply, err = r.do("HSET", redisInt, key, fieldk, fieldv)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) HGet(key, field string) (res string, err error) {
	var reply interface{}
	reply, err = r.do("HGET", redisString, key, field)
	if err != nil {
		return
	}
	res = reply.(string)
	return
}

func (r *Redis) HGetInt(key, field string) (res int, err error) {
	var reply interface{}
	reply, err = r.do("HGET", redisInt, key, field)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) HMGet(key string, fields ...interface{}) (res []string, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, fields...)
	reply, err = r.do("HMGET", redisStrings, keys...)
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

func (r *Redis) HMSet(key string, fields ...interface{}) (res string, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, fields...)
	reply, err = r.do("HMSET", redisString, keys...)
	if err != nil {
		return
	}
	res = reply.(string)
	return
}

func (r *Redis) HGetAll(key string) (res map[string]string, err error) {
	var reply interface{}
	reply, err = r.do("HGETALL", redisStringMap, key)
	if err != nil {
		return
	}
	res = reply.(map[string]string)
	return
}

func (r *Redis) HKeys(key string) (res []string, err error) {
	var reply interface{}
	reply, err = r.do("HKEYS", redisStrings, key)
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

/*
*	set
 */
func (r *Redis) SAdd(key string, members ...interface{}) (res int, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, members...)
	reply, err = r.do("SADD", redisInt, keys...)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) SRem(key string, members ...interface{}) (res int, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, members...)
	reply, err = r.do("SREM", redisInt, keys...)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) SIsMember(key string, member string) (res bool, err error) {
	var reply interface{}
	reply, err = r.do("SISMEMBER", redisBool, key, member)
	if err != nil {
		return
	}
	res = reply.(bool)

	return
}

func (r *Redis) SMembers(key string) (res []string, err error) {
	var reply interface{}
	reply, err = r.do("SMEMBERS", redisStrings, key)
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

/*
	ZSET
*/
func (r *Redis) ZAdd(key string, args ...interface{}) (res int, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, args...)
	reply, err = r.do("ZADD", redisInt, keys...)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) ZRange(key string, args ...interface{}) (res []string, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, args...)
	reply, err = r.do("ZRANGE", redisStrings, keys...)
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

func (r *Redis) ZRangeInt(key string, start, stop int) (res []int, err error) {
	var reply interface{}
	reply, err = r.do("ZRANGE", redisInts, key, start, stop)
	if err != nil {
		return
	}
	res = reply.([]int)
	return
}

func (r *Redis) ZRangeWithScore(key string, start, stop int) (res []string, err error) {
	var reply interface{}
	reply, err = r.do("ZRANGE", redisStrings, key, start, stop, "WITHSCORES")
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

func (r *Redis) ZCount(key string, min, max int) (res int, err error) {
	var reply interface{}
	reply, err = r.do("ZCOUNT", redisInt, key, min, max)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) ZCard(key string) (res int, err error) {
	var reply interface{}
	reply, err = r.do("ZCARD", redisInt, key)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) LLen(key string) (res int64, err error) {
	var reply interface{}
	reply, err = r.do("LLEN", redisInt64, key)
	if err != nil {
		return
	}
	res = reply.(int64)
	return
}

func (r *Redis) Incrby(key string, incr int) (res int64, err error) {
	var reply interface{}
	reply, err = r.do("INCRBY", redisInt64, key, incr)
	if err != nil {
		return
	}
	res = reply.(int64)
	return
}

func (r *Redis) ZIncrby(key string, incr int, member string) (res int, err error) {
	var reply interface{}
	reply, err = r.do("ZINCRBY", redisInt, key, incr, member)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

/*
* If the member not in the zset or key not exits, ZRank will return ErrNil
 */
func (r *Redis) ZRank(key string, member string) (res int, err error) {
	var reply interface{}
	reply, err = r.do("ZRANK", redisInt, key, member)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

/*
* If the members not in the zset or key not exits, ZRem will return ErrNil
 */
func (r *Redis) ZRem(key string, members ...interface{}) (res int, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, members...)

	reply, err = r.do("ZREM", redisInt, keys...)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) ZRemrangebyrank(key string, members ...interface{}) (res int, err error) {
	var reply interface{}
	keys := []interface{}{key}
	keys = append(keys, members...)

	reply, err = r.do("ZREMRANGEBYRANK", redisInt, keys...)
	if err != nil {
		return
	}
	res = reply.(int)
	return
}

func (r *Redis) Subscribe(ctx context.Context, key string, maxSize int) (chan []byte, error) {
	ch := make(chan []byte, maxSize)

	if r.opts.ReadTimeout < 100 && r.opts.ReadTimeout > 0 {
		return ch, errors.New("Read timeout should be longer")
	}

	healthCheckPeriod := r.opts.ReadTimeout * 70 / 100

	var offHealthCheck = (healthCheckPeriod == 0)
	done := make(chan error, 1)

	// While not a permanent error on the connection.
	go func() {
	start:
		client := r.pool.Get()
		// defer client.Close()
		psc := redis.PubSubConn{client}
		// Set up subscriptions
		err := psc.Subscribe(key)
		if err != nil {
			return
		}

		go func(c redis.PubSubConn) {
			if offHealthCheck {
				return
			}
			ticker := time.NewTicker(time.Duration(healthCheckPeriod * 10e5))
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err = c.Ping(""); err != nil {
						break
					}
				case <-ctx.Done():
					return
				case <-done:
					return
				}
			}
		}(psc)

		for client.Err() == nil {
			select {
			case <-ctx.Done():
				client.Close()
				return
			default:
				switch v := psc.Receive().(type) {
				case redis.Message:
					ch <- v.Data
				case redis.Subscription:
					log.Infof("Receive chan (%s) %s %d", v.Channel, v.Kind, v.Count)
				case error:
					log.Errorf("Receive error (%v), client will reconnect..", v)
					client.Close()
					if !offHealthCheck {
						done <- v
					}
					time.Sleep(time.Second / 10)
					goto start
				}
			}
		}
	}()

	return ch, nil
}

/*
* If the member not in the zset or key not exits, ZScore will return ErrNil
 */
func (r *Redis) ZScore(key, member string) (res float64, err error) {
	var reply interface{}
	reply, err = r.do("ZSCORE", redisFloat64, key, member)
	if err != nil {
		return
	}
	res = reply.(float64)
	return
}

func (r *Redis) Zrevrange(key string, args ...interface{}) (res []string, err error) {
	var reply interface{}
	argss := []interface{}{key}
	argss = append(argss, args...)
	reply, err = r.do("ZREVRANGE", redisStrings, argss...)
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

func (r *Redis) Zrevrangebyscore(key string, args ...interface{}) (res []string, err error) {
	var reply interface{}
	argss := []interface{}{key}
	argss = append(argss, args...)
	reply, err = r.do("ZREVRANGEBYSCORE", redisStrings, argss...)
	if err != nil {
		return
	}
	res = reply.([]string)
	return
}

func (r *Redis) ZrevrangebyscoreInt(key string, args ...interface{}) (res []int, err error) {
	var reply interface{}
	argss := []interface{}{key}
	argss = append(argss, args...)
	reply, err = r.do("ZREVRANGEBYSCORE", redisInts, argss...)
	if err != nil {
		return
	}
	res = reply.([]int)
	return
}

// Pipelining

func (r *Redis) randomDuration(n int64) time.Duration {
	s := rand.NewSource(r.lastTime)
	return time.Duration(rand.New(s).Int63n(n) + 1)
}

func (r *Redis) do(cmd string, f func(interface{}, error) (interface{}, error), args ...interface{}) (reply interface{}, err error) {
	stCode := redisSuccess
	defer func() {
		atomic.StoreInt64(&r.lastTime, time.Now().UnixNano())
	}()
	count := 0
	now := time.Now()

retry1:
	client := r.pool.Get()
	defer client.Close()
	if client.Err() == redis.ErrPoolExhausted {
		if r.opts.Retry > 0 && count < r.opts.Retry {
			count = count + 1
			goto retry1
		}
		stCode = redisConnExhausted
		return nil, ErrConnExhausted
	}
	if err = client.Err(); err != nil {
		stCode = redisConnError
		if r.opts.Retry > 0 && count < r.opts.Retry {
			count = count + 1
			time.Sleep(time.Millisecond * r.randomDuration(10))
			goto retry1
		}
		return nil, err
	}

retry2:
	reply, err = client.Do(cmd, args...)
	if r.opts.Retry > 0 && count < r.opts.Retry && err != nil && err != redis.ErrNil {
		count = count + 1
		log.GenLogf("redisclient retry %d times, cmd %s cause %s", count, cmd, err)
		time.Sleep(time.Millisecond * r.randomDuration(10))
		goto retry2
	}

	if f != nil {
		reply, err = f(reply, err)
	}
	address := r.opts.Addr

	if err == redis.ErrNil {
		err = nil
	} else if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			stCode = redisTimeout
			err = ErrTimeout
		} else {
			stCode = redisError
		}
		log.GenLogf("%d|redisclient|%s|%s|%d|%s|%s", serverLocalPid, "reqid", cmd, stCode, err, address)
	}

	endTime := time.Now()
	costTime := time.Now().Sub(now).Nanoseconds() / 1e6
	if (r.opts.SlowTime > 0 && costTime > int64(r.opts.SlowTime)) || (stCode == redisTimeout) {
		log.SlowLogf("%d|%s|redisclient|%s|%s|%d|%d|%s|%s", serverLocalPid, endTime.Format(logFormat), "reqid", cmd, stCode, costTime, address, "nil")
	}
	return
}
