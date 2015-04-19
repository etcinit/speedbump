package speedbump

import (
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
