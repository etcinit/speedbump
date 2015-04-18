package speedbump

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gopkg.in/redis.v2"
)

func createClient() *redis.Client {
	return redis.NewTCPClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func Test_NewLimiter(t *testing.T) {
	client := createClient()
	hasher := PerSecondHasher{}

	NewLimiter(client, hasher, 10)
}

func Test_Attempts(t *testing.T) {
	client := createClient()
	hasher := PerSecondHasher{}

	limiter := NewLimiter(client, hasher, 5)

	has, err := limiter.Has("127.0.0.1")
	assert.Nil(t, err)
	assert.False(t, has)

	ok, err := limiter.Attempt("127.0.0.1")
	assert.Nil(t, err)
	assert.True(t, ok)

	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.2")
	limiter.Attempt("127.0.0.1")
	limiter.Attempt("127.0.0.2")
	ok, err = limiter.Attempt("127.0.0.1")

	assert.Nil(t, err)
	assert.False(t, ok)

	left, err := limiter.Left("127.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, 0, left)

	left, err = limiter.Left("127.0.0.2")
	assert.Nil(t, err)
	assert.True(t, left > 0)

	ok, err = limiter.Attempt("127.0.0.2")
	assert.Nil(t, err)
	assert.True(t, ok)

	time.Sleep(time.Second)

	left, err = limiter.Left("127.0.0.1")
	assert.Nil(t, err)
	assert.Equal(t, 5, left)
}
