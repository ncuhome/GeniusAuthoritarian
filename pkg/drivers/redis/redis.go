package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func New(conf *Config) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:       conf.Addr,
		Password:   conf.Password,
		DB:         conf.DB,
		MaxConnAge: time.Hour * 5,
	})
	return c, c.Ping(context.Background()).Err()
}
