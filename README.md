# What is this?

This is a configurable rate limiter middleware for Golang.

It is also one of the final challenges of https://goexpert.fullcycle.com.br/pos-goexpert/.

# How it works?

In its recommended configuration (using Redis Storage Adapter) it uses a Redis sorted set as a sliding window to keep record of an IP or token last accesses. If more than the allowed requests per second are done, a block record is created, and the user must wait until the block time ends.

# How to test?

Run it with docker compose:

```bash
docker compose up -d   
```

If debug mode is active, you can follow the logs with:

```bash
docker compose logs server --follow 
```

You can manually request at `http://localhost:8080` or use a stress tester, like fortio. Some useful fortio command lines for the bundled ".env":

```bash

# get a block for IP access
docker run  \
--rm \
--network=rate-limiter \
fortio/fortio load -logger-force-color -qps 501 -c 1 -t 1s http://rate-limiter-server:8080

# get a block for generic token access
docker run  \
--rm \
--network=rate-limiter \
fortio/fortio load -logger-force-color -qps 2000 -c 1 -t 510ms -H "API_KEY: zzz" http://rate-limiter-server:8080

# get a block for "abc" token access
docker run  \
--rm \
--network=rate-limiter \
fortio/fortio load -logger-force-color -qps 11 -c 1 -t 1s -H "API_KEY: abc" http://rate-limiter-server:8080

# get a block for "def" token access
docker run  \
--rm \
--network=rate-limiter \
fortio/fortio load -logger-force-color -qps 601 -c 1 -t 1s -H "API_KEY: def" http://rate-limiter-server:8080

# smash the webserver
docker run  \
--rm \
--network=rate-limiter \
fortio/fortio load -logger-force-color -qps 10000 -c 100 -t 1s http://rate-limiter-server:8080

```

# How to configure?

## Environment Variables

The easiest way to configure is through environment variables. You can set them in your OS, in docker-compose.yml or just editing the `.env` file:

|Value|Type|Description|Default Value|
|---|---|---|---|
|RATE_LIMITER_IP_MAX_REQUESTS|integer|Requests per second allowed for an IP.|100|
|RATE_LIMITER_IP_BLOCK_TIME|integer|Block time in milliseconds for IPs that reach their request quota.|1000|
|RATE_LIMITER_TOKEN_MAX_REQUESTS|integer|Requests per second allowed for a token (any token). This has priority over IP configuration.|200|
|RATE_LIMITER_TOKEN_BLOCK_TIME|integer|Block time in milliseconds for tokens (any token) that reach their request quota. This has priority over IP configuration.|500|
|RATE_LIMITER_TOKEN_AAA_MAX_REQUESTS|integer|Requests per second allowed for the token "AAA". This has priority over token configuration. If not defined, it will use RATE_LIMITER_TOKEN_MAX_REQUESTS for this token. |-|
|RATE_LIMITER_TOKEN_AAA_BLOCK_TIME|integer|Block time in milliseconds for the token "AAA" when it reachs its request quota. This has priority over token configuration. If not defined, it will use RATE_LIMITER_TOKEN_BLOCK_TIME for this token. |-|
|RATE_LIMITER_DEBUG|boolean|Runs in debug mode. A lot of messages are displayed on stdout.|false|
|RATE_LIMITER_USE_REDIS|boolean|Uses the Redis Storage Adapter.|false|
|RATE_LIMITER_REDIS_ADDRESS|string|Redis host for Redis Storage Adapter.|-|
|RATE_LIMITER_REDIS_PASSWORD|string|Redis password for Redis Storage Adapter.|-|
|RATE_LIMITER_REDIS_DB|integer|Redis database for Redis Storage Adapter.|-|

With environment variables there is no need to pass anything directly to the middleware. Just create it:

```go
rateLimiter := ratelimiter.NewRateLimiter()
```

## Code Configuration

You can use code configuration if you want. Environment variables will override code values, but you can disable this behavior setting `DisableEnvs: true`:

```go
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

```

## Custom Adapters

You can write a custom Storage Adapter (store accesses and blocks) and Response Writer (write the status codes and messages to the request).

You can use `./ratelimiter/adapter/redis_storage_adapter.go` and `ratelimiter/responsewriter/default_response_writer.go` as base to write yours. You can set them with code configuration:

```
rateLimiter := ratelimiter.NewRateLimiterWithConfig(
	&ratelimiter.RateLimiterConfig{
		StorageAdapter: myCustomStorageAdapter{},
		ResponseWriter: myCustomResponseWriter{},
	},
)

```

# How to use?

Whatever way you configure your middleware, you use it as any other midlleware. Example with go-chi:

```go
rateLimiter := ratelimiter.NewRateLimiter()
r := chi.NewRouter()
r.Use(rateLimiter)
```
