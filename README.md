# [speedbump](https://github.com/etcinit/speedbump)

A Redis-backed Rate Limiter for Go

[![wercker status](https://app.wercker.com/status/9832225d9e89d9702d4ce7ca4e8e4285/m "wercker status")](https://app.wercker.com/project/bykey/9832225d9e89d9702d4ce7ca4e8e4285)

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
	"gopkg.in/redis.v2"
)

func main() {
	client := redis.NewTCPClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	hasher := speedbump.PerSecondHasher{}

	// Here we create a limiter that will only allow 10 requests per second
	limiter := speedbump.NewLimiter(client, hasher, 5)

	for {
		// This example has a hardcoded IP, but you would replace it with the IP
		// of a client on a real case.
		success, err := limiter.Attempt("127.0.0.1")

		if err != nil {
			panic(err)
		}

		//fmt.Println(left)
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
