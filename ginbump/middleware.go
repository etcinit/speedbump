package ginbump

import (
	"net"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/etcinit/speedbump"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v2"
)

// RateLimit is a Gin middleware for rate limitting incoming requests based on
// the client's IP address.
func RateLimit(client *redis.Client, hasher speedbump.RateHasher, max int64) gin.HandlerFunc {
	limiter := speedbump.NewLimiter(client, hasher, max)

	return func(c *gin.Context) {
		// Attempt to perform the request
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		ok, err := limiter.Attempt(ip)

		if err != nil {
			panic(err)
		}

		if !ok {
			nextTime := time.Now().Add(hasher.Duration())

			c.JSON(429, gin.H{
				"status":   "error",
				"messages": []string{"Rate limit exceeded. Try again in " + humanize.Time(nextTime)},
			})
			c.Abort()
		}

		c.Next()

		// After the request
		// log.Print(ip + " was limited because it exceeded the max rate")
	}
}
