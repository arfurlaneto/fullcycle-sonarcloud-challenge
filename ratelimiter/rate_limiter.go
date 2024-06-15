package ratelimiter

import (
	"context"
	"time"
)

func checkRateLimit(ctx context.Context, keyType string, key string, config *RateLimiterConfig, rateConfig *RateLimiterRateConfig) (*time.Time, error) {
	if key == "" {
		return nil, nil
	}

	block, err := config.StorageAdapter.GetBlock(ctx, keyType, key)
	if err != nil {
		return nil, err
	}

	if block == nil {
		success, count, err := config.StorageAdapter.IncrementAccesses(ctx, keyType, key, rateConfig.MaxRequestsPerSecond)
		if err != nil {
			return nil, err
		}

		if success {
			DebugPrintf(config, "%d of %d (%dms if blocked)", keyType, key, count, rateConfig.MaxRequestsPerSecond, rateConfig.BlockTimeMilliseconds)
		} else {
			DebugPrintf(config, "adding a block of %dms", keyType, key, rateConfig.BlockTimeMilliseconds)
			block, err = config.StorageAdapter.AddBlock(ctx, keyType, key, rateConfig.BlockTimeMilliseconds)
			if err != nil {
				return nil, err
			}
		}
	}

	if block != nil {
		DebugPrintf(config, "block time %.2f seconds", keyType, key, GetRemainingBlockTime(block))
		return block, nil
	}

	return nil, nil
}
