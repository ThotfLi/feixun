package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	RedisURL            = "localhost:6379"
	redisMaxIdle        = 16   //最大空闲连接数
	redisIdleTimeoutSec = 240 //最大空闲连接时间
	RedisPassword       = ""
)

var RedisPool *redis.Pool

func NewRedisPool() *redis.Pool{
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",RedisURL)
		},
	}
}






















