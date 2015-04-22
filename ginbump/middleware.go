// Package ginbump provides an example Speedbump middleware for the Gin
// framework.
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
//
// The resulting middleware will use the client to talk to the Redis server.
// The hasher is used to keep track of counters and to provide an estimate of
// when the client should be able to do requests again. The limit per period is
// defined by the max.
//
// Response format
//
// Once a client reaches the imposed limit, they will receive a JSON response
// similar to the following:
//
//  {
//    "messages":["Rate limit exceeded. Try again in 1 minute from now"],
//    "status":"error"
//  }
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
