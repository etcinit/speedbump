// Package speedbump provides a Redis-backed rate limiter.
package speedbump

import (
	"strconv"
	"time"

	"gopkg.in/redis.v5"
)

// RateLimiter is a Redis-backed rate limiter.
type RateLimiter struct {
	// redisClient is the client that will be used to talk to the Redis server.
	redisClient *redis.Client
	// hasher is used to generate keys for each counter and to set their
	// expiration time.
	hasher RateHasher
	// max defines the maximum number of attempts that can occur during a
	// period.
	max int64
}

// RateHasher is an object capable of generating a hash that uniquely identifies
// a counter to track the number of requests for an id over a certain time
// interval. The input of the Hash function can be any unique id, such as an IP
// address.
type RateHasher interface {
	// Hash is the hashing function.
	Hash(id string) string
	// Duration returns the duration of each period. This is used to determine
	// when to expire each counter key, and can also be used by other libraries
	// to generate messages that provide an estimate of when the limit will
	// expire.
	Duration() time.Duration
}

// NewLimiter creates a new instance of a rate limiter.
func NewLimiter(client *redis.Client, hasher RateHasher, max int64) *RateLimiter {
	return &RateLimiter{
		redisClient: client,
		hasher:      hasher,
		max:         max,
	}
}

// Has returns whether the rate limiter has seen a request for a specific id
// during the current period.
func (r *RateLimiter) Has(id string) (bool, error) {
	hash := r.hasher.Hash(id)
	return r.redisClient.Exists(hash).Result()
}

// Attempted returns the number of attempted requests for an id in the current
// period. Attempted does not count attempts that exceed the max requests in an
// interval and only returns the max count after this is reached.
func (r *RateLimiter) Attempted(id string) (int64, error) {
	hash := r.hasher.Hash(id)
	val, err := r.redisClient.Get(hash).Result()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist. See: http://redis.io/commands/GET
			return 0, nil
		}
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(val, 10, 64)
}

// Left returns the number of remaining requests for id during a current period.
func (r *RateLimiter) Left(id string) (int64, error) {
	// Retrieve attempted count.
	attempted, err := r.Attempted(id)
	if err != nil {
		return 0, err
	}
	// Left is max minus attempted.
	left := r.max - attempted
	if left < 0 {
		return 0, nil
	}
	return left, nil
}

// Attempt attempts to perform a request for an id and returns whether it was
// successful or not.
func (r *RateLimiter) Attempt(id string) (bool, error) {
	// Create hash from id
	hash := r.hasher.Hash(id)
	// Get value for hash in Redis. If redis.Nil is returned, key does not exist.
	exists := true
	val, err := r.redisClient.Get(hash).Result()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist. See: http://redis.io/commands/GET
			exists = false
		} else {
			return false, err
		}
	}
	// If key exists and is >= max requests, return false.
	if exists {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return false, err
		}
		if intVal >= r.max {
			return false, nil
		}
	}
	// Otherwise, increment and expire key for hasher.Duration(). Note, we call
	// Expire even when key already exists to avoid race condition where key
	// expires between prior existence check and this Incr call.
	// See: http://redis.io/commands/INCR
	// See: http://redis.io/commands/INCR#pattern-rate-limiter-1
	err = r.redisClient.Watch(func(rx *redis.Tx) error {
		_, err := rx.Pipelined(func(pipe *redis.Pipeline) error {
			if err := pipe.Incr(hash).Err(); err != nil {
				return err
			}
			if err := pipe.Expire(hash, r.hasher.Duration()).Err(); err != nil {
				return err
			}
			return nil
		})
		return err
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
