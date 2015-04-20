# ginbump

Example Speedbump middleware for Gin

## Usage:

Somewhere in your Gin engine setup code:

```go

// Create a Redis client
client := redis.NewTCPClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Limit the engine's or group's requests to a maximum of 100 requests per
// client per minute.
engineOrGroup.Use(ginbump.RateLimit(client, speedbump.PerMinuteHasher{}, 100))
```

after that, if clients stay within the limit, the won't notice anything. If they
do go over the limit, the will get an HTTP 429 error (Too Many Requests) with
the following content:

```js
{
    "messages":["Rate limit exceeded. Try again in 1 minute from now"],
    "status":"error"
}
```
