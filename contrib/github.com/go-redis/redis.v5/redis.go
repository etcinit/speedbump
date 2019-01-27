package redis

import (
	"time"

	"github.com/ntindall/speedbump/internal"
	redis "gopkg.in/redis.v5"
)

// Wrapper is a wrapper around *redis.Client that implements the
// internal.RedisClient interface.
type Wrapper struct {
	*redis.Client
}

var _ internal.RedisClient = &Wrapper{}

func (w *Wrapper) Exists(key string) (exists bool, err error) {
	return w.Client.Exists(key).Result()
}

func (w *Wrapper) Get(key string) (value string, err error) {
	return w.Client.Get(key).Result()
}

func (w *Wrapper) IncrAndExpire(key string, duration time.Duration) error {
	return w.Client.Watch(func(rx *redis.Tx) error {
		_, err := rx.Pipelined(func(pipe *redis.Pipeline) error {
			if err := pipe.Incr(key).Err(); err != nil {
				return err
			}

			return pipe.Expire(key, duration).Err()
		})

		return err
	})
}

// NewRedisClient constructs a speedbump.RedisClient from a "gopkg.in/redis.v5"
// redis.Client.
func NewRedisClient(redisClient *redis.Client) internal.RedisClient {
	return &Wrapper{redisClient}
}
