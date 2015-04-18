package speedbump

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_PerSecond_Hash(t *testing.T) {
	hasher := PerSecondHasher{}

	resultOne := hasher.Hash("127.0.0.1")

	time.Sleep(time.Second)

	resultTwo := hasher.Hash("127.0.0.1")

	assert.NotEqual(t, resultOne, resultTwo)

	resultThree := hasher.Hash("127.0.0.1")
	resultFour := hasher.Hash("127.0.0.1")
	resultFive := hasher.Hash("127.0.0.2")

	assert.Equal(t, resultThree, resultFour)
	assert.NotEqual(t, resultFour, resultFive)
}

func Test_PerSecond_Duration(t *testing.T) {
	hasher := PerSecondHasher{}

	assert.Equal(t, time.Second, hasher.Duration())
}
