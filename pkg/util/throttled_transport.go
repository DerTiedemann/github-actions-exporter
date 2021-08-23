package util

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type ThrottledTransport struct {
	roundTripperWrap http.RoundTripper
	ratelimiter      *rate.Limiter
}

func (c *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	err := c.ratelimiter.Wait(r.Context()) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.roundTripperWrap.RoundTrip(r)
	remaining, ok := resp.Header["X-Ratelimit-Remaining"]
	if ok {
		b, err := strconv.ParseInt(remaining[0], 10, 64)
		if err != nil {
			log.Println(err)
		} else {
			c.ratelimiter.SetBurst(int(b))
		}
	}

	return resp, err
}

// NewThrottledTransport wraps transportWrap with a rate limitter
// examle usage:
// client := http.DefaultClient
// client.Transport = NewThrottledTransport(10*time.Seconds, 60, http.DefaultTransport) allows 60 requests every 10 seconds
func NewThrottledTransport(limitPeriod time.Duration, requestCount int, transportWrap http.RoundTripper) http.RoundTripper {
	return &ThrottledTransport{
		roundTripperWrap: transportWrap,
		ratelimiter:      rate.NewLimiter(rate.Every(limitPeriod), requestCount),
	}
}
