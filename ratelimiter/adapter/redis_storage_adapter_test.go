package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RateLimitRedisStorageAdapter struct {
	suite.Suite
	context context.Context
}

func TestRateLimitRedisStorageAdapter(t *testing.T) {
	suite.Run(t, new(RateLimitRedisStorageAdapter))
}

func (s *RateLimitRedisStorageAdapter) SetupTest() {
	s.context = context.Background()
}

func (s *RateLimitRedisStorageAdapter) TestNewRateLimitRedisStorageAdapter() {
	storageAdapter := NewRateLimitRedisStorageAdapter("", "", 0)
	assert.NotNil(s.T(), storageAdapter)
}

func (s *RateLimitRedisStorageAdapter) TestFormatRedisKey() {
	storageAdapter := NewRateLimitRedisStorageAdapter("", "", 0)
	redisKeys := storageAdapter.formatRedisKey("block", "uSeR-ToKeN", "AbC123*#")
	assert.Equal(s.T(), "block-user_token-AbC123*#", redisKeys)
}
