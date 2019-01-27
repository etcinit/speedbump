package redis

import (
	"time"

	redis "github.com/gomodule/redigo/redis"
	"github.com/ntindall/speedbump/internal"
)

type redisWrapper struct {
	conn redis.Conn
}

var _ internal.RedisClient = &redisWrapper{}

func (w *redisWrapper) Exists(key string) (exists bool, err error) {
	return redis.Bool(w.conn.Do("EXISTS", key))
}

func (w *redisWrapper) Get(key string) (value string, err error) {
	return redis.String(w.conn.Do("GET", key))
}

func (w *redisWrapper) IncrAndExpire(key string, duration time.Duration) error {
	if err := w.conn.Send("MULTI"); err != nil {
		return err
	}
	if err := w.conn.Send("INCR", key); err != nil {
		return err
	}
	if err := w.conn.Send("EXPIRE", key, duration/time.Second); err != nil {
		return err
	}
	_, err := w.conn.Do("EXEC")
	return err
}

// NewRedisClient constructs a internal.RedisClient from a redigo connection.
func NewRedisClient(redisConn redis.Conn) internal.RedisClient {
	return &redisWrapper{conn: redisConn}
}
