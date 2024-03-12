package adapter

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type rateLimitRedisStorageAdapter struct {
	client *redis.Client
}

func NewRateLimitRedisStorageAdapter(address string, password string, db int64) *rateLimitRedisStorageAdapter {
	adapter := rateLimitRedisStorageAdapter{}

	adapter.client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       int(db),
	})

	return &adapter
}

func (s *rateLimitRedisStorageAdapter) IncrementAccesses(ctx context.Context, keyType string, key string, maxAccesses int64) (bool, int64, error) {
	redisKey := s.formatRedisKey("access", keyType, key)

	now := time.Now()
	clearBefore := now.Add(-time.Second)

	pipeline := s.client.Pipeline()

	pipeline.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(clearBefore.UnixMicro(), 10))
	count := pipeline.ZCard(ctx, redisKey)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		logRedisError(err)
		return false, 0, err
	}

	if count.Val() >= maxAccesses {
		return false, count.Val(), nil
	}

	pipeline = s.client.Pipeline()

	pipeline.ZAdd(ctx, redisKey, redis.Z{Score: float64(now.UnixMicro()), Member: now.Format(time.RFC3339Nano)})
	pipeline.Expire(ctx, redisKey, time.Second)

	_, err = pipeline.Exec(ctx)
	if err != nil {
		logRedisError(err)
		return false, 0, err
	}

	return true, count.Val() + 1, nil
}

func (s *rateLimitRedisStorageAdapter) GetBlock(ctx context.Context, keyType string, key string) (*time.Time, error) {
	redisKey := s.formatRedisKey("block", keyType, key)

	value, err := s.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		logRedisError(err)
		return nil, err
	}

	parsedValue, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return nil, err
	}

	return &parsedValue, nil
}

func (s *rateLimitRedisStorageAdapter) AddBlock(ctx context.Context, keyType string, key string, milliseconds int64) (*time.Time, error) {
	redisKey := s.formatRedisKey("block", keyType, key)

	expiration := time.Duration(int64(time.Millisecond) * milliseconds)
	blockedUntil := time.Now().Add(expiration)

	_, err := s.client.Set(ctx, redisKey, blockedUntil.Format(time.RFC3339Nano), expiration).Result()
	if err != nil {
		logRedisError(err)
		return nil, err
	}

	return &blockedUntil, nil
}

func (s *rateLimitRedisStorageAdapter) formatRedisKey(prefix string, keyType string, key string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		strings.ToLower(prefix),
		strings.ToLower(strings.ReplaceAll(keyType, "-", "_")),
		key,
	)
}

func logRedisError(err error) {
	fmt.Printf(
		"%s [REDIS STORAGE ADAPTER] ERROR: %s\n",
		time.Now().UTC().Format("2006-01-02 15:04:05"),
		err.Error(),
	)
}
