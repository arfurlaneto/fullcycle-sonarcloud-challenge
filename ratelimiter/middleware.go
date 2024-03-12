package ratelimiter

import (
	"context"
	"net"
	"net/http"
	"time"
)

type rateLimiterCheckFunction = func(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error)

func NewRateLimiter() func(next http.Handler) http.Handler {
	return NewRateLimiterWithConfig(nil)
}

func NewRateLimiterWithConfig(config *RateLimiterConfig) func(next http.Handler) http.Handler {
	config = setConfiguration(config)
	return func(next http.Handler) http.Handler {
		return rateLimiter(config, next, checkRateLimit)
	}
}

func rateLimiter(config *RateLimiterConfig, next http.Handler, checkRateLimitFn rateLimiterCheckFunction) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var block *time.Time
		var err error

		token := r.Header.Get("API_KEY")
		if token != "" {
			tokenConfig, _ := config.GetRateLimiterRateConfigForToken(token)
			block, err = checkRateLimitFn(r.Context(), "TOKEN", token, config, tokenConfig)
		} else {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			block, err = checkRateLimitFn(r.Context(), "IP", host, config, config.IP)
		}

		if err != nil {
			config.ResponseWriter.WriteError(&w, err)
			return
		}

		if block != nil {
			config.ResponseWriter.WriteResponse(&w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
