package ginbump

import (
	"net/http"

	"github.com/etcinit/speedbump"
	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v5"
)

// The following example shows how to set up a rate limitting middleware in Gin
// that allows 100 requests per client per minute.
func ExampleRateLimit() {
	// Create a Gin engine
	router := gin.Default()

	// Add a route
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Limit the engine's requests to a maximum of 100 requests per client per
	// minute.
	router.Use(RateLimit(client, speedbump.PerMinuteHasher{}, 100))

	// Start listening
	router.Run(":8080")
}
