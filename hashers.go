package speedbump

import (
	"strconv"
	"time"
)

// PerSecondHasher generates hashes per second. This means you can keep track
// of N request per second.
type PerSecondHasher struct{}

// Hash generates the hash for the current period and client.
func (h PerSecondHasher) Hash(id string) string {
	return id + ":" + strconv.FormatInt(time.Now().Unix(), 10)
}

// Duration gets the duration of each period.
func (h PerSecondHasher) Duration() time.Duration {
	return time.Second
}
