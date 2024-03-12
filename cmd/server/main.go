package main

import (
	"net/http"

	"github.com/arfurlaneto/fullcycle-sonarcloud-challenge/ratelimiter"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	rateLimiter := ratelimiter.NewRateLimiterWithConfig(
		&ratelimiter.RateLimiterConfig{
			IP: &ratelimiter.RateLimiterRateConfig{
				MaxRequestsPerSecond:  100,  // same as RATE_LIMITER_IP_MAX_REQUESTS
				BlockTimeMilliseconds: 5000, // same as RATE_LIMITER_IP_BLOCK_TIME
			},
			Token: &ratelimiter.RateLimiterRateConfig{
				MaxRequestsPerSecond:  500, // same as RATE_LIMITER_TOKEN_MAX_REQUESTS
				BlockTimeMilliseconds: 500, // same as RATE_LIMITER_TOKEN_BLOCK_TIME
			},
			// same as RATE_LIMITER_TOKEN_AAA_MAX_REQUESTS and RATE_LIMITER_TOKEN_AAA_BLOCK_TIME
			CustomTokens: &map[string]*ratelimiter.RateLimiterRateConfig{
				"ABC_1": {MaxRequestsPerSecond: 2000, BlockTimeMilliseconds: 100},
				"ABC_2": {MaxRequestsPerSecond: 2000, BlockTimeMilliseconds: 100},
			},
			Debug:       true, // same as RATE_LIMITER_DEBUG
			DisableEnvs: true, // if true, environment values are ignored
		},
	)

	r := chi.NewRouter()

	r.Use(rateLimiter)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
