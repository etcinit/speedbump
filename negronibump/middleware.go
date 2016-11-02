package negronibump

import (
	"net"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dustin/go-humanize"
	"github.com/etcinit/speedbump"
	"github.com/unrolled/render"
	"gopkg.in/redis.v5"
)

func RateLimit(client *redis.Client, hasher speedbump.RateHasher, max int64) negroni.HandlerFunc {
	limiter := speedbump.NewLimiter(client, hasher, max)
	rnd := render.New()

	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		ok, err := limiter.Attempt(ip)
		if err != nil {
			panic(err)
		}

		if !ok {
			nextTime := time.Now().Add(hasher.Duration())
			rnd.JSON(rw, 429, map[string]string{"error": "Rate limit exceeded. Try again in " + humanize.Time(nextTime)})
		} else {
			next(rw, r)
		}
	}
}
