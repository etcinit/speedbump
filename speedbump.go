// Package speedbump provides a Redis-backed rate limiter.
package speedbump

import (
	"strconv"
	"time"

	"gopkg.in/redis.v2"
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

// RateHasher is an object capable of generating a hash that uniquely
// identifies a counter that keeps track of the number of requests attempted by
// a client on a period of time. The input of the function can be anything that
// can uniquely identify a client, but it usually an IP address.
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

// Has returns whether the rate limiter has seen/received a request from a
// specific client during the current period.
func (r *RateLimiter) Has(id string) (bool, error) {
	hash := r.hasher.Hash(id)

	return r.redisClient.Exists(hash).Result()
}

// Attempted returns the number of attempted requests for a client in the
// current period.
//
// Not all attempts will be recorded, once the limit has been reached, the
// counter will stop adding up.
func (r *RateLimiter) Attempted(id string) (int64, error) {
	has, err := r.Has(id)

	if err != nil {
		return 0, err
	}

	if !has {
		return 0, nil
	}

	hash := r.hasher.Hash(id)
	str, err := r.redisClient.Get(hash).Result()

	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(str, 10, 64)
}

// Left returns the number of remaining requests for client during a current
// period.
func (r *RateLimiter) Left(id string) (int64, error) {
	attempted, err := r.Attempted(id)

	if err != nil {
		return 0, nil
	}

	left := r.max - attempted

	if left < 0 {
		return 0, nil
	}

	return left, nil
}

// Attempt attempts to perform a request for a client and returns whether it
// was successful or not.
func (r *RateLimiter) Attempt(id string) (bool, error) {
	hash := r.hasher.Hash(id)

	exists, err := r.Has(id)

	if err != nil {
		return false, err
	}

	if exists {
		str, err := r.redisClient.Get(hash).Result()

		if err != nil {
			return false, err
		}

		intVal, err := strconv.ParseInt(str, 10, 64)

		if err != nil {
			return false, err
		}

		if str != "" && intVal > r.max {
			return false, nil
		}

		err = r.redisClient.Incr(hash).Err()

		if err != nil {
			return false, err
		}

		return true, nil
	}

	rx := r.redisClient.Multi()
	defer rx.Close()

	_, err = rx.Exec(func() error {
		if err := rx.Incr(hash).Err(); err != nil {
			return err
		}

		if err := rx.Expire(hash, r.hasher.Duration()).Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
