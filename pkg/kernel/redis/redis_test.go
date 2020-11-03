package redis

import (
	"testing"

	"context"
	"log"
	"time"
)

var (
	store *Redis
)

func init() {
	var err error
	store, err = NewRedis(&RedisConfig{
		Addr:        "localhost:6379",
		ServerName:  "haha",
		ReadTimeout: 1000,
	})

	if err != nil {
		log.Fatalf("err:%s", err)
	}
}

func TestPipelining(t *testing.T) {
	pip, err := store.NewPipelining()
	if err != nil {
		t.Fatalf("err:%s", err)
	}
	if err := pip.Send("SET", 1, 1); err != nil {
		t.Fatalf("err:%s", err)
	}
	if err := pip.Send("GET", 1); err != nil {
		t.Fatalf("err:%s", err)
	}
	if err := pip.Flush(); err != nil {
		t.Fatalf("err:%s", err)
	}

	r1, err := String(pip.Receive())
	if err != nil {
		t.Fatalf("err:%s", err)
	} else if r1 != "OK" {
		t.Fatalf("err:%s", "r1 != OK")
	}
	r2, err := String(pip.Receive())
	if err != nil {
		t.Fatalf("err:%s", err)
	} else if r2 != "1" {
		t.Fatalf("err:%s", "r2 != 1")
	}
	_, err = String(pip.Receive())
	if err == nil {
		t.Fatalf("err:%s", "err == nil")
	}
	pip.Close()
}

func TestSrem(t *testing.T) {
	a := []interface{}{"BLACK_USER-28", 3217428, 124324}
	reply, err := store.Do("srem", a...)
	log.Printf("?? = %v %v\n", reply, err)
}

func TestSend(t *testing.T) {
	store.Del("TestSend")
	if err := store.Send("TestSend", "this is value"); err != nil {
		t.Fatalf("err:%s", err)
	}

	store.Del("TestSend")
}

func TestReceive(t *testing.T) {
	store.Del("TestReceive")
	for i := 0; i < 10; i++ {
		store.Send("TestReceive", i)
	}
	c := make(chan struct{})
	receiveCh := store.Receive("TestReceive", c, 10)
	i := 0
	go func() {
		defer close(c)
		time.Sleep(time.Second * 2)
	}()
	for b := range receiveCh {
		if b == nil {
			t.Fatalf("result not match")
		}
		i++
	}
	if i != 10 {
		t.Fatalf("result not match i=%d", i)
	}
	store.Del("TestReceive")
}

func TestSet(t *testing.T) {
	if ret, err := store.Set("TestSet", "this is value"); err != nil {
		t.Fatalf("err:%s", err)
	} else if !ret {
		t.Fatalf("ret not match")
	}

	store.Del("TestSet")
}

func TestSetExSecond(t *testing.T) {
	if ret, err := store.SetExSecond("TestSetExSecond", "this is value", 2); err != nil {
		t.Fatalf("err:%s", err)
	} else if ret != "OK" {
		t.Fatalf("ret not match")
	}

	time.Sleep(2 * time.Second)

	if ret, err := store.Get("TestSetExSecond"); err != nil {
		t.Fatalf("err=%s\n", err)
	} else if string(ret) != "" {
		t.Fatalf("ret not match")
	}
	store.Del("TestSetExSecond")
}

func TestGet(t *testing.T) {
	store.Set("TestGet", "this is value")
	if ret, err := store.Get("TestGet"); err != nil {
		t.Fatalf("err:%s", err)
	} else if string(ret) != "this is value" {
		t.Fatalf("ret not match")
	}

	store.Set("TestGet", 888)
	if ret, err := store.GetInt("TestGet"); err != nil {
		t.Fatalf("err:%s", err)
	} else if ret != 888 {
		t.Fatalf("ret not match")
	}

	//	store.Set("NOTEXIST", "")
	if ret, err := store.Get("NOTEXIST"); err != nil {
		t.Fatalf("err:%s", err)
	} else if string(ret) != "" {
		t.Fatalf("err")
	}
	store.Del("TestGet")
	//	store.Del("NOTEXIST")
}

func TestMGet(t *testing.T) {
	store.Set("TestMGet1", "this is value1")
	store.Set("TestMGet2", "this is value2")
	if ret, err := store.MGet("TestMGet1", "TestMGet2"); err != nil {
		t.Fatalf("err:%s", err)
	} else if string(ret[0]) != "this is value1" ||
		string(ret[1]) != "this is value2" {
		t.Fatalf("ret not match")
	}
	if ret, err := store.MGet("TestMGetNotExist", "TestMGet2NOTEXIST"); err != nil {
		t.Fatalf("err:%s", err)
	} else if len(ret) != 2 {
		t.Fatalf("ret not match")
	}
	//store.Del("TestMGet1")
	//store.Del("TestMGet2")
}

func TestMSet(t *testing.T) {
	if ret, err := store.MSet("TestMSet1", 1, "TestMSet2", 2); err != nil {
		t.Fatalf("err:%s", err)
	} else if ret != "OK" || err != nil {
		t.Fatalf("ret not match")
	}
	store.Del("TestMSet1")
	store.Del("TestMSet2")
}

func TestDel(t *testing.T) {
	store.Set("TestDel1", "this is value1")
	store.Set("TestDel2", "this is value2")
	if count, err := store.Del("TestDel1", "TestDel2"); err != nil {
		t.Fatalf("err:%s count:%d\n", err, count)
	} else if count != 2 {
		t.Fatalf("count not match")
	}
}

func TestExist(t *testing.T) {
	if exist, err := store.Exists("TestExist"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if exist {
		t.Fatalf("not match")
	}
	store.Set("TestExist", "this is value")
	if exist, err := store.Exists("TestExist"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if !exist {
		t.Fatalf("not match")
	}
	store.Del("TestExist")
}

func TestExpire(t *testing.T) {
	store.Set("TestExpire", "this is value2")
	if err := store.Expire("TestExpire", time.Second*3); err != nil {
		t.Fatalf("err:%s\n", err)
	}

	if s, err := store.Get("TestExpire"); err != nil {
		t.Fatalf("%s\n", err)
	} else if string(s) != "this is value2" {
		t.Fatalf("not match")
	}

	time.Sleep(time.Second * 3)
	if _, err := store.Get("TestExpire"); err != nil {
		t.Fatalf("%s\n", err)
	}
}

func TestHGet(t *testing.T) {
	store.HSet("TestHGet", "KEY1", "V1")
	if res, err := store.HGet("TestHGet", "KEY1"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != "V1" {
		t.Fatalf("not match")
	}

	store.HSet("TestHGet", "KEY2", 123)
	if res, err := store.HGetInt("TestHGet", "KEY2"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != 123 {
		t.Fatalf("not match")
	}
	store.Del("TestHGet")
}

func TestHDel(t *testing.T) {
	store.HSet("TestHDel", "KEY2", "V2")
	if res, err := store.HDel("TestHDel", "KEY2"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != 1 {
		t.Fatalf("not match")
	}

	if res, err := store.HDel("TestHDel", "KEY1"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != 0 {
		t.Fatalf("not match")
	}
	store.Del("TestHDel")
}

func TestHKeys(t *testing.T) {
	store.HSet("TestHkeys", "k1", "v1")
	if ret, err := store.HKeys("TestHkeys"); err != nil {
		t.Fatalf("err:%s", err)
	} else {
		t.Logf("%v\n", ret)
	}
}

func TestHSet(t *testing.T) {
	store.Del("TestHSet")
	if res, err := store.HSet("TestHSet", "KEY1", "V1"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != 1 {
		t.Fatalf("err: res not 1, res=%d\n", res)
	}

	if res, err := store.HSet("TestHSet", "KEY2", 1234); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != 1 {
		t.Fatalf("err: res not 1, res=%d\n", res)
	}
	store.Del("TestHSet")
}

func TestHMGet(t *testing.T) {
	if res, err := store.HMGet("TestHMGet", "KEY3", "KEY4", "KEY5"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if len(res) != 3 {
		t.Fatalf("err: res not 3, res=%v len=%d\n", res, len(res))
	}
	store.HSet("TestHMGet", "KEY3", "V3")
	store.HSet("TestHMGet", "KEY4", "V4")
	if res, err := store.HMGet("TestHMGet", "KEY3", "KEY4"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if len(res) != 2 {
		t.Fatalf("err: res not 2, res=%v\n", res)
	} else if res[0] != "V3" || res[1] != "V4" {
		t.Fatalf("result not match")
	}
	store.Del("TestHMGet")
}

func TestHMSet(t *testing.T) {
	if res, err := store.HMSet("TestHMSet", "KEY3", "KEY4"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if res != "OK" {
		t.Fatalf("err: res not 3, res=%v len=%d\n", res, len(res))
	}
	store.Del("TestHMSet")
}

func TestHGetAll(t *testing.T) {
	store.Del("TestHGetAll")
	if res, err := store.HGetAll("TestHGetAll"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else if len(res) != 0 {
		t.Fatalf("err: res not 0, res=%v len=%d\n", res, len(res))
	}
	store.HSet("TestHGetAll", "KEY3", "V3")
	store.HSet("TestHGetAll", "KEY4", "V4")
	store.HSet("TestHGetAll", "KEY5", "V5")
	if m, err := store.HGetAll("TestHGetAll"); err != nil {
		t.Fatalf("err:%s\n", err)
	} else {
		if m["KEY3"] != "V3" || m["KEY4"] != "V4" || m["KEY5"] != "V5" {
			t.Fatalf("result not match\n")
		}
	}
	store.Del("TestHGetAll")
}

func TestSAdd(t *testing.T) {
	store.Del("TestSAdd")
	store.SAdd("TestSAdd", "M1")
	if res, err := store.SAdd("TestSAdd", "M1", "M2", "M3"); err != nil {
		t.Fatalf("TestSAdd: %s\n", err)
	} else if res != 2 {
		t.Fatalf("result not match res=%d\n", res)
	}
	store.Del("TestSAdd")
}

func TestSRem(t *testing.T) {
	store.Del("TestSRem")
	store.SAdd("TestSRem", "M1")
	if res, err := store.SRem("TestSRem", "M1"); err != nil {
		t.Fatalf("TestSRem: %s\n", err)
	} else if res != 1 {
		t.Fatalf("result not match res=%d\n", res)
	}
	store.Del("TestSRem")
}

func TestSIsMember(t *testing.T) {
	store.Del("TestSIsMember")
	store.SAdd("TestSIsMember", "M1")
	if res, err := store.SIsMember("TestSIsMember", "M1"); err != nil {
		t.Fatalf("TestSIsMember %s\n", err)
	} else if !res {
		t.Fatalf("result not match res=%v\n", res)
	}
	if res, err := store.SIsMember("TestSIsMember", "M2"); err != nil {
		t.Fatalf("TestSIsMember %s\n", err)
	} else if res {
		t.Fatalf("result not match res=%v\n", res)
	}
	store.Del("TestSIsMember")
}

func TestSMembers(t *testing.T) {
	store.Del("TestSMembers")
	store.SAdd("TestSMembers", "m1", "m2", "m3")
	if res, err := store.SMembers("TestSMembers"); err != nil {
		t.Fatalf("TestSMembers %s\n", err)
	} else if len(res) != 3 {
		t.Fatalf("result not match res=%v\n", res)
	}
	store.Del("TestSMembers")
	if res, err := store.SMembers("TestSMembers"); err != nil {
		t.Fatalf("TestSMembers %s\n", err)
	} else if len(res) != 0 {
		t.Fatalf("result not match res=%v\n", res)
	}

}

func TestZSet(t *testing.T) {
	store.Del("TestZSet")
	store.ZAdd("TestZSet", 2, "m1")
	if res, err := store.ZAdd("TestZSet", 1, "m1", 2, "m2"); err != nil {
		t.Fatalf("TestZSet %s\n", err)
	} else if res != 1 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if _, err := store.ZAdd("TestZSet", 1, "m1", 2, "m2", 4); err == nil {
		t.Fatalf("TestZSet err\n")
	}
	store.Del("TestZSet")
}

func TestZRange(t *testing.T) {
	store.Del("TestZRange")

	store.ZAdd("TestZRange", 1, "m1", 1.1, "m11", 2, "m2", 3, "m3", 4, "m4")
	if res, err := store.ZRange("TestZRange", 0, 3); err != nil {
		t.Fatalf("TestZRange %s\n", err)
	} else if len(res) != 4 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if res, err := store.ZRangeWithScore("TestZRange", 0, 3); err != nil {
		t.Fatalf("TestZRange %s\n", err)
	} else if len(res) != 8 {
		t.Fatalf("result not match res=%v\n", res)
	}

	store.Del("TestZRange")
}

func TestZCount(t *testing.T) {
	store.Del("TestZCount")

	store.ZAdd("TestZCount", 1, "m1", 2, "m2", 3, "m3", 4, "m4")
	if res, err := store.ZCount("TestZCount", 0, 3); err != nil {
		t.Fatalf("TestZCount %s\n", err)
	} else if res != 3 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if res, err := store.ZCount("TestZRange", 10, 3); err != nil {
		t.Fatalf("TestZCount %s\n", err)
	} else if res != 0 {
		t.Fatalf("result not match res=%v\n", res)
	}

	store.Del("TestZCount")
}

func TestZCard(t *testing.T) {
	store.Del("TestZCard")

	store.ZAdd("TestZCard", 1, "m1", 2, "m2", 3, "m3", 4, "m4")
	if res, err := store.ZCard("TestZCard"); err != nil {
		t.Fatalf("TestZCard %s\n", err)
	} else if res != 4 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if res, err := store.ZCard("TestZCardNone"); err != nil {
		t.Fatalf("TestZCard %s\n", err)
	} else if res != 0 {
		t.Fatalf("result not match res=%v\n", res)
	}
	store.Del("TestZCard")
}

func TestZIncrby(t *testing.T) {
	store.Del("TestZIncrby")

	store.ZAdd("TestZIncrby", 1, "m1")
	if res, err := store.ZIncrby("TestZIncrby", 10, "m1"); err != nil {
		t.Fatalf("TestZIncrby %s\n", err)
	} else if res != 11 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if res, err := store.ZIncrby("TestZIncrby", 10, "m2"); err != nil {
		t.Fatalf("TestZIncrby %s\n", err)
	} else if res != 10 {
		t.Fatalf("result not match res=%v\n", res)
	}

	store.Del("TestZIncrby")
}

func TestZRank(t *testing.T) {
	store.Del("TestZRank")

	store.ZAdd("TestZRank", 1, "m1", 2, "m2")
	if res, err := store.ZRank("TestZRank", "m2"); err != nil {
		t.Fatalf("TestZRank %s\n", err)
	} else if res != 1 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if res, err := store.ZRank("TestZRank", "m3"); err != nil {
		t.Fatalf("result not match res=%v\n", res)
	}

	if res, err := store.ZRank("TestZRankxxxx", "m3"); err != nil {
		t.Fatalf("result not match res=%v\n", res)
	}

	store.Del("TestZRank")
}

func TestZRem(t *testing.T) {
	store.Del("TestZRem")

	store.ZAdd("TestZRem", 1, "m1", 2, "m2", 3, "m3")
	if res, err := store.ZRem("TestZRem", "m2", "m3", "m4"); err != nil {
		t.Fatalf("TestZRem %s\n", err)
	} else if res != 2 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if _, err := store.ZRank("TestZRemxxx", "m3"); err != nil {
		t.Fatalf("result not match err=%s\n", err)
	}

	if _, err := store.ZRank("TestZRem", "m33"); err != nil {
		t.Fatalf("TestZRem err: %v\n", err)
	}

	store.Del("TestZRem")
}

func TestZScore(t *testing.T) {
	store.Del("TestZScore")

	store.ZAdd("TestZScore", 1, "m1", 2, "m2", 3.3, "m3")
	if res, err := store.ZScore("TestZScore", "m3"); err != nil {
		t.Fatalf("TestZRem %s\n", err)
	} else if res != 3.3 {
		t.Fatalf("result not match res=%v\n", res)
	}

	if _, err := store.ZScore("TestZScore", "m4"); err != nil {
		t.Fatalf("TestZScore err: %v\n", err)
	}

	store.Del("TestZScore")
}

func TestTimeOut(t *testing.T) {
	timeout := 200 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	reply, err := store.DoCtx(ctx, "BLPOP", "3214321", 1)
	if err != context.DeadlineExceeded {
		t.Fatalf("TestTimeOut err: %s %v\n", err, reply)
	}

	timeout = 5 * time.Second
	ctx, _ = context.WithTimeout(context.Background(), timeout)
	_, err = store.DoCtx(ctx, "BLPOP", "3214321", 0)
	if err != ErrTimeout {
		t.Fatalf("TestTimeOut err: %s\n", err)
	}
}

func TestZrevrangebyscoreInt(t *testing.T) {
	store.Del("uids")
	store.ZAdd("uids", 1, "1")
	store.ZAdd("uids", 2, "2")
	store.ZAdd("uids", 3, "3")
	store.ZAdd("uids", 4, "4")
	store.ZAdd("uids", 5, "5")
	store.ZAdd("uids", 6, "6")

	uids, err := store.ZrevrangebyscoreInt("uids", "+inf", "-inf", "LIMIT", 0, 3)
	if err != nil {
		t.Fatal(err)
	}
	if uids[0] != 6 || uids[1] != 5 || uids[2] != 4 {
		t.Fatal(uids)
	}
	store.Del("uids")
}
