package speedbump

import (
	"fmt"
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/stretchr/testify/assert"
)

func Test_PerSecond_Hash(t *testing.T) {
	mock := clock.NewMock()
	hasher := PerSecondHasher{
		Clock: mock,
	}

	resultOne := hasher.Hash("127.0.0.1")

	mock.Add(time.Second)

	resultTwo := hasher.Hash("127.0.0.1")

	assert.NotEqual(t, resultOne, resultTwo)

	resultThree := hasher.Hash("127.0.0.1")
	resultFour := hasher.Hash("127.0.0.1")
	resultFive := hasher.Hash("127.0.0.2")

	assert.Equal(t, resultThree, resultFour)
	assert.NotEqual(t, resultFour, resultFive)

	// Test that it can create a new clock
	hasher = PerSecondHasher{}
	hasher.Hash("127.0.0.1")
}

// The following example shows how to create mock hashers for testing the rate
// limiter in your code:
func ExamplePerSecondHasher() {
	// Create a mock clock.
	mock := clock.NewMock()

	// Create a new per second hasher with the mock clock.
	hasher := PerMinuteHasher{
		Clock: mock,
	}

	// Generate two consecutive hashes. On most systems, the following should
	// generate two identical hashes.
	hashOne := hasher.Hash("127.0.0.1")
	hashTwo := hasher.Hash("127.0.0.1")

	// Now we push the clock forward by a minute (time travel).
	mock.Add(time.Minute)

	// The third hash should be different now.
	hashThree := hasher.Hash("127.0.0.1")

	fmt.Println(hashOne == hashTwo)
	fmt.Println(hashOne == hashThree)
	// Output: true
	// false
}

func Test_PerSecond_Duration(t *testing.T) {
	hasher := PerSecondHasher{}

	assert.Equal(t, time.Second, hasher.Duration())
}

func Test_PerMinute_Hash(t *testing.T) {
	mock := clock.NewMock()
	hasher := PerMinuteHasher{
		Clock: mock,
	}

	resultOne := hasher.Hash("127.0.0.1")

	mock.Add(time.Minute)

	resultTwo := hasher.Hash("127.0.0.1")

	assert.NotEqual(t, resultOne, resultTwo)

	resultThree := hasher.Hash("127.0.0.1")
	resultFour := hasher.Hash("127.0.0.1")
	resultFive := hasher.Hash("127.0.0.2")

	assert.Equal(t, resultThree, resultFour)
	assert.NotEqual(t, resultFour, resultFive)

	// Test that it can create a new clock
	hasher = PerMinuteHasher{}
	hasher.Hash("127.0.0.1")
}

func Test_PerMinute_Duration(t *testing.T) {
	hasher := PerMinuteHasher{}

	assert.Equal(t, time.Minute, hasher.Duration())
}

func Test_PerHour_Hash(t *testing.T) {
	mock := clock.NewMock()
	hasher := PerHourHasher{
		Clock: mock,
	}

	resultOne := hasher.Hash("127.0.0.1")

	mock.Add(time.Hour)

	resultTwo := hasher.Hash("127.0.0.1")

	assert.NotEqual(t, resultOne, resultTwo)

	resultThree := hasher.Hash("127.0.0.1")
	resultFour := hasher.Hash("127.0.0.1")
	resultFive := hasher.Hash("127.0.0.2")

	assert.Equal(t, resultThree, resultFour)
	assert.NotEqual(t, resultFour, resultFive)

	// Test that it can create a new clock
	hasher = PerHourHasher{}
	hasher.Hash("127.0.0.1")
}

func Test_PerHour_Duration(t *testing.T) {
	hasher := PerHourHasher{}

	assert.Equal(t, time.Hour, hasher.Duration())
}
