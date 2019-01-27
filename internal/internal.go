package internal

import (
	"time"
)

// RedisClient is an abstraction over speedbump connection to redis.
// It is exported from internal so that it can only be instructed from
// within the package.
type RedisClient interface {
	Get(key string) (value string, err error)
	Exists(key string) (exists bool, err error)
	IncrAndExpire(key string, duration time.Duration) error
}
