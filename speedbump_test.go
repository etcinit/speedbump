package speedbump

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/redis.v5"
)

func createClient() *redis.Client {
	if os.Getenv("WERCKER_REDIS_HOST") != "" {
		return redis.NewClient(&redis.Options{
			Addr:     os.Getenv("WERCKER_REDIS_HOST") + ":6379",
			Password: "",
			DB:       0,
		})
	}

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func teardown(t *testing.T, client *redis.Client) {
	// Flush Redis.
	require.NoError(t, client.FlushAll().Err())
}

func TestNewLimiter(t *testing.T) {
	client := createClient()
	hasher := PerSecondHasher{}
	max := int64(10)
	actual := NewLimiter(client, hasher, max)

	assert.Exactly(t, RateLimiter{
		redisClient: client,
		hasher:      hasher,
		max:         max,
	}, *actual)
}

func ExampleNewLimiter() {
	// Create a Redis client.
	client := createClient()

	// Create a new hasher.
	hasher := PerSecondHasher{}

	// Create a new limiter that will only allow 10 requests per second.
	limiter := NewLimiter(client, hasher, 10)

	fmt.Println(limiter.Attempt("127.0.0.1"))
	// Output: true <nil>
}

func TestHas(t *testing.T) {
	// Create Redis client and defer DB teardown.
	client := createClient()
	defer teardown(t, client)
	// Create limiter of 5 requests/min.
	limiter := NewLimiter(client, PerMinuteHasher{}, 5)
	// Choose an arbitrary id.
	testID := "test_id"

	// Ensure testID returns false initially.
	has, err := limiter.Has(testID)
	require.NoError(t, err)
	assert.False(t, has)

	// Make request attempts, including attempt for testID.
	_, err = limiter.Attempt("some_id")
	require.NoError(t, err)
	_, err = limiter.Attempt(testID)
	require.NoError(t, err)
	_, err = limiter.Attempt("some_other_id")
	require.NoError(t, err)

	// Ensure testID returns true.
	has, err = limiter.Has(testID)
	require.NoError(t, err)
	assert.True(t, has)

	// Ensure other ids return false.
	has, err = limiter.Has("some_unseen_id")
	require.NoError(t, err)
	assert.False(t, has)
}

type attempt struct {
	id          string
	expectedOK  bool
	expectedErr string
}

func TestAttempt(t *testing.T) {
	// Create Redis client and defer DB teardown.
	client := createClient()
	defer teardown(t, client)
	// Create PerMinuteHasher with mock clock.
	mock := clock.NewMock()
	hasher := PerMinuteHasher{
		Clock: mock,
	}
	// Create limiter of 5 requests/min.
	limiter := NewLimiter(client, hasher, 5)
	// Choose an arbitrary id.
	testID := "test_id"
	// Ensure no key exists before first request for testID.
	has, err := limiter.Has(testID)
	require.NoError(t, err)
	assert.False(t, has)

	// Set up series of requests that include > max requests for testID.
	attempts := []attempt{
		// Attempt #1 for testID.
		{
			id:          testID,
			expectedOK:  true,
			expectedErr: "",
		},
		// Attempt #2 for testID.
		{
			id:          testID,
			expectedOK:  true,
			expectedErr: "",
		},
		{
			id:          "some_id",
			expectedOK:  true,
			expectedErr: "",
		},
		{
			id:          "some_id",
			expectedOK:  true,
			expectedErr: "",
		},
		// Attempt #3 for testID.
		{
			id:          testID,
			expectedOK:  true,
			expectedErr: "",
		},
		// Attempt #4 for testID.
		{
			id:          testID,
			expectedOK:  true,
			expectedErr: "",
		},
		{
			id:          "some_other_id",
			expectedOK:  true,
			expectedErr: "",
		},
		// Attempt #5 for testID.
		{
			id:          testID,
			expectedOK:  true,
			expectedErr: "",
		},
		// Attempt #6 for testID -- exceeded max.
		{
			id:          testID,
			expectedOK:  false,
			expectedErr: "",
		},
		// Attempt #7 for testID -- exceeded max.
		{
			id:          testID,
			expectedOK:  false,
			expectedErr: "",
		},
		{
			id:          "some_id",
			expectedOK:  true,
			expectedErr: "",
		},
		// Attempt #8 for testID -- exceeded max.
		{
			id:          testID,
			expectedOK:  false,
			expectedErr: "",
		},
	}

	// Make attempts and keep track of results. Avoid asserting in loop to ensure
	// execution takes < 1 minute.
	resultsOK := []bool{}
	resultsErr := []error{}
	for i := 0; i < len(attempts); i++ {
		ok, err := limiter.Attempt(attempts[i].id)
		resultsOK = append(resultsOK, ok)
		resultsErr = append(resultsErr, err)
	}
	// Assert attempt results.
	for i := 0; i < len(resultsErr); i++ {
		if attempts[i].expectedErr == "" {
			require.NoError(t, resultsErr[i], "failed case #%d: %+v", i, attempts[i])
			assert.Exactly(t, attempts[i].expectedOK, resultsOK[i], "failed case #%d: %+v", i, attempts[i])
		} else {
			require.Error(t, resultsErr[i], "failed case #%d: %+v", i, attempts[i])
			assert.Contains(t, resultsErr[i].Error(), attempts[i].expectedErr, "failed case #%d: %+v", i, attempts[i])
		}
	}

	// Mock add 1 minute to simulate waiting 1 minute, expect true for testID.
	mock.Add(time.Minute)
	ok, err := limiter.Attempt(attempts[0].id)
	require.NoError(t, err)
	assert.True(t, ok, "Attempts returned false after waiting for interval")
}

func makeNAttempts(t *testing.T, limiter *RateLimiter, id string, n int64) {
	var i int64
	for i = 0; i < n; i++ {
		_, err := limiter.Attempt(id)
		require.NoError(t, err, "got error during request attempt")
	}
}

func TestAttemptedLeft(t *testing.T) {
	// Create Redis client and defer DB teardown.
	client := createClient()
	defer teardown(t, client)
	// Create PerMinuteHasher with mock clock.
	mock := clock.NewMock()
	hasher := PerMinuteHasher{
		Clock: mock,
	}
	max := int64(5)
	// Create limiter of 5 requests/min.
	limiter := NewLimiter(client, hasher, max)
	// Choose an arbitrary id.
	testID := "test_id"

	// Check we have max left, 0 attempted initially.
	{
		left, err := limiter.Left(testID)
		require.NoError(t, err)
		assert.Exactly(t, max, left)
		attempted, err := limiter.Attempted(testID)
		require.NoError(t, err)
		assert.Exactly(t, int64(0), attempted)
	}

	// Make 2 attempts.
	makeNAttempts(t, limiter, testID, 2)
	{
		// Expect max-2 left for testID.
		left, err := limiter.Left(testID)
		require.NoError(t, err)
		assert.Exactly(t, max-2, left)
		// Expect max left for any other id.
		leftOther, err := limiter.Left("some_id")
		require.NoError(t, err)
		assert.Exactly(t, max, leftOther)
		// Expect 2 attempted for testID.
		attempted, err := limiter.Attempted(testID)
		require.NoError(t, err)
		assert.Exactly(t, int64(2), attempted)
		// Expect 0 attempted for any other id.
		attemptedOther, err := limiter.Attempted("some_id")
		require.NoError(t, err)
		assert.Exactly(t, int64(0), attemptedOther)
	}

	// Make max-2 more attempts.
	makeNAttempts(t, limiter, testID, max-2)
	{
		// Expect 0 left for testID.
		left, err := limiter.Left(testID)
		require.NoError(t, err)
		assert.Exactly(t, int64(0), left)
		// Expect max left for any other id.
		leftOther, err := limiter.Left("some_id")
		require.NoError(t, err)
		assert.Exactly(t, max, leftOther)
		// Expect max attempted for testID.
		attempted, err := limiter.Attempted(testID)
		require.NoError(t, err)
		assert.Exactly(t, max, attempted)
		// Expect 0 attempted for any other id.
		attemptedOther, err := limiter.Attempted("some_id")
		require.NoError(t, err)
		assert.Exactly(t, int64(0), attemptedOther)

	}

	// Make 10 more attempts.
	makeNAttempts(t, limiter, testID, 10)
	{
		// Expect 0 left for testID.
		left, err := limiter.Left(testID)
		require.NoError(t, err)
		assert.Exactly(t, int64(0), left)
		// Expect max left for any other id.
		leftOther, err := limiter.Left("some_id")
		require.NoError(t, err)
		assert.Exactly(t, max, leftOther)
		// Expect max attempted for testID.
		attempted, err := limiter.Attempted(testID)
		require.NoError(t, err)
		assert.Exactly(t, max, attempted)
		// Expect 0 attempted for any other id.
		attemptedOther, err := limiter.Attempted("some_id")
		require.NoError(t, err)
		assert.Exactly(t, int64(0), attemptedOther)
	}

	// Mock add 1 minute to simulate waiting 1 minute.
	mock.Add(time.Minute)
	{
		// Expect max left for testID.
		left, err := limiter.Left(testID)
		require.NoError(t, err)
		assert.Exactly(t, max, left)
		// Expect max left for any other id.
		leftOther, err := limiter.Left("some_id")
		require.NoError(t, err)
		assert.Exactly(t, max, leftOther)
		// Expect 0 attempted for testID.
		attempted, err := limiter.Attempted(testID)
		require.NoError(t, err)
		assert.Exactly(t, int64(0), attempted)
		// Expect 0 attempted for any other id.
		attemptedOther, err := limiter.Attempted("some_id")
		require.NoError(t, err)
		assert.Exactly(t, int64(0), attemptedOther)

	}
}
