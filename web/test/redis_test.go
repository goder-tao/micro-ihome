package test

import (
	"context"
	"ihome/web/dao/redis"
	"log"
	"testing"
	"time"
	r 	"github.com/go-redis/redis/v8"
)

func TestRedis(t *testing.T)  {
	redis.Init()
	defer redis.Close()
	cli, err := redis.Get()
	ctx := context.Background()
	cli.Set(ctx, "test", "ok", 0)
	v, err := cli.Get(ctx, "test").Result()

	// test simple set
	if err != nil {
		t.Error("Get: ", err)
	}
	if v != "ok" {
		t.Error("value is not equal")
	}

	// test expiration
	cli.SetEX(ctx, "ex", "ok", time.Second*2)

	v, err = cli.Get(ctx, "ex").Result()
	if err != nil {
		t.Error("Ex Get: ", err)
	}
	if v != "ok" {
		t.Error("ex value is not equal")
	}

	time.Sleep(time.Second*3)
	v, err = cli.Get(ctx, "ex").Result()
	if err != nil && err == r.Nil {
		log.Print("Ex correct")
	}
}
