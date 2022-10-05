package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var cli *redis.Client
var mu sync.Mutex
var one sync.Once

const (
	KEY_AREA              = "AREA"
	KEY_FACILITY          = "facilities"
	KEY_PREFIX_USERAVATAR = "user_avatar"

	EXPIRATION_USERAVATAR = time.Minute
)

func Init() error {
	var err error = nil
	one.Do(func() {
		cli = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Username: "",
			Password: "",
			DB:       0,
			PoolSize: 100,
		})
		ctx := context.Background()
		cli.Set(ctx, "inti", "ok", time.Second)
		_, e := cli.Get(ctx, "init").Result()
		if e != nil {
			err = e
		}
	})
	return err
}

func Get() (*redis.Client, error) {
	if cli == nil {
		mu.Lock()
		defer mu.Unlock()
		if cli == nil {
			err := Init()
			if err != nil {
				return nil, err
			}
		}
		return cli, nil
	}
	return cli, nil
}

func Close() {
	if cli != nil {
		cli.Close()
	}
}
