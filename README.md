# [speedbump](https://github.com/etcinit/speedbump) [![GoDoc](https://godoc.org/github.com/etcinit/speedbump?status.svg)](http://godoc.org/github.com/etcinit/speedbump)

A Redis-backed Rate Limiter for Go

[![wercker status](https://app.wercker.com/status/9832225d9e89d9702d4ce7ca4e8e4285/m/master "wercker status")](https://app.wercker.com/project/bykey/9832225d9e89d9702d4ce7ca4e8e4285)

## Cool stuff

- Backed by Redis, so it keeps track of requests across a cluster
- Extensible timing functions. Includes defaults for tracking requests per
second, minute, and hour
- Works with IPv4, IPv6, or any other unique identifier
- Example middleware included for [Gin](https://github.com/gin-gonic/gin) (See: [ginbump](https://github.com/etcinit/speedbump/blob/master/ginbump)) and
[Negroni](https://github.com/codegangsta/negroni) (See:
[negronibump](https://github.com/etcinit/speedbump/blob/master/negronibump))

## Notes

- This library is fairly new and could break
- The current implementation could have key leakage/race conditions (See: http://redis.io/commands/incr#pattern-rate-limiter)

## Usage

- Get a working Redis server
- Go get:

```sh
$ go get github.com/etcinit/speedbump
```

- Include it in your code

```go
package main

import (
	"fmt"
	"time"

	"github.com/etcinit/speedbump"
	"gopkg.in/redis.v5"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	hasher := speedbump.PerSecondHasher{}

	// Here we create a limiter that will only allow 5 requests per second
	limiter := speedbump.NewLimiter(client, hasher, 5)

	for {
		// This example has a hardcoded IP, but you would replace it with the IP
		// of a client on a real case.
		success, err := limiter.Attempt("127.0.0.1")

		if err != nil {
			panic(err)
		}

		if success {
			fmt.Println("Successful!")
		} else {
			fmt.Println("Limited! :(")
		}

		time.Sleep(time.Millisecond * time.Duration(100))
	}
}
```

- Output:

```
Successful!
Successful!
Successful!
Successful!
Successful!
Successful!
Limited! :(
Limited! :(
Limited! :(
Limited! :(
Limited! :(
Successful!
Successful!
Successful!
Successful!
Successful!
Successful!
Limited! :(
Limited! :(
Limited! :(
Limited! :(
Successful!
Successful!
...
```
